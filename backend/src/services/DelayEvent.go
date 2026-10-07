package services

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"groundTurn/src/constants"
	"groundTurn/src/constructors"
	"groundTurn/src/middlewares"
	"groundTurn/src/repositories"
	"groundTurn/src/types"
	"groundTurn/src/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DelayEventService 登记、关闭、归因延误。
// 顺延规则：登记的延误默认未关闭，立即计入该航班任务的 effective_deadline；
// 关闭（resolved_at 置值）后这几分钟不再顺延。全程不改任务 deadline 列。
type DelayEventService struct {
	delayRepo      *repositories.DelayEventRepository
	turnaroundRepo *repositories.FlightTurnaroundRepository
}

func NewDelayEventService(dr *repositories.DelayEventRepository, tr *repositories.FlightTurnaroundRepository) *DelayEventService {
	return &DelayEventService{delayRepo: dr, turnaroundRepo: tr}
}

// List 返回延误视图，附每个航班当前未关闭累计分钟。
func (s *DelayEventService) List() ([]types.DelayEventView, error) {
	events, err := s.delayRepo.List()
	if err != nil {
		return nil, err
	}
	ids := []uint{}
	seen := map[uint]bool{}
	for i := range events {
		if !seen[events[i].TurnaroundID] {
			seen[events[i].TurnaroundID] = true
			ids = append(ids, events[i].TurnaroundID)
		}
	}
	delayMap, err := s.delayRepo.MapByTurnaround(ids)
	if err != nil {
		return nil, err
	}
	summaryByTAR := map[uint]utils.OpenDelaySummary{}
	for _, id := range ids {
		summaryByTAR[id] = utils.SummarizeOpenDelays(delayMap[id])
	}
	views := make([]types.DelayEventView, 0, len(events))
	for i := range events {
		flightNo := ""
		if f, err := s.turnaroundRepo.Get(events[i].TurnaroundID); err == nil {
			flightNo = f.FlightNo
		}
		views = append(views, constructors.BuildDelayEventView(events[i], flightNo, summaryByTAR[events[i].TurnaroundID]))
	}
	return views, nil
}

// Register 登记一笔未关闭延误，分钟数立即顺延该航班任务的生效截止。
func (s *DelayEventService) Register(c *gin.Context, req types.DelayEventRegisterRequest) (*types.DelayEventView, error) {
	if !slices.Contains(constants.DelayType, req.DelayType) {
		return nil, utils.EnumError("delay_type", req.DelayType)
	}
	if req.Minutes <= 0 {
		return nil, utils.ValidationError(constants.DeadlineNotExtendedMsg)
	}
	if _, err := s.turnaroundRepo.Get(req.TurnaroundID); err != nil {
		return nil, utils.NotFoundError(constants.TargetTurnaround, req.TurnaroundID)
	}
	now := time.Now().UTC()
	event := constructors.NewDelayEvent(req, now)
	if err := s.delayRepo.Create(&event); err != nil {
		return nil, err
	}

	// 登记延误不改写任何任务 deadline；这里仅汇总当前未关闭分钟供日志/响应展示。
	summary, err := s.currentSummary(req.TurnaroundID)
	if err != nil {
		return nil, err
	}
	flightNo := ""
	if f, ferr := s.turnaroundRepo.Get(req.TurnaroundID); ferr == nil {
		flightNo = f.FlightNo
	}
	view := constructors.BuildDelayEventView(event, flightNo, summary)

	middlewares.QueueAudit(c, "DELAY_REGISTER", constants.TargetDelay, fmt.Sprint(event.ID),
		utils.RenderLog(constants.LogTemplates["DELAY_REGISTER"], map[string]interface{}{
			"TurnaroundID": req.TurnaroundID, "DelayType": req.DelayType, "Minutes": req.Minutes,
			"OpenDelayMinutes": summary.OpenMinutes, "ResponsibilityTeam": req.ResponsibilityTeam,
		}))
	middlewares.QueueAudit(c, "DELAY_DEADLINE_RECALC", constants.TargetDelay, fmt.Sprint(event.ID),
		utils.RenderLog(constants.LogTemplates["DELAY_DEADLINE_RECALC"], map[string]interface{}{
			"TurnaroundID": req.TurnaroundID, "OpenCount": summary.OpenCount,
			"OpenDelayMinutes": summary.OpenMinutes,
		}))
	return &view, nil
}

// Resolve 关闭延误：resolved_at 置当前时间，该笔分钟从顺顺延口径中剔除。
// 已关闭的延误重复关闭返回 409。
func (s *DelayEventService) Resolve(c *gin.Context, id uint) (*types.DelayEventView, error) {
	event, err := s.delayRepo.Get(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, utils.NotFoundError(constants.TargetDelay, id)
	}
	if err != nil {
		return nil, err
	}
	if event.ResolvedAt != nil {
		return nil, utils.ValidationError(fmt.Sprintf("delay #%d already resolved at %s", id, event.ResolvedAt.Format(time.RFC3339)))
	}
	now := time.Now().UTC()
	if err := s.delayRepo.Resolve(id, now); err != nil {
		return nil, err
	}
	event.ResolvedAt = &now
	summary, err := s.currentSummary(event.TurnaroundID)
	if err != nil {
		return nil, err
	}
	flightNo := ""
	if f, ferr := s.turnaroundRepo.Get(event.TurnaroundID); ferr == nil {
		flightNo = f.FlightNo
	}
	view := constructors.BuildDelayEventView(*event, flightNo, summary)

	middlewares.QueueAudit(c, "DELAY_RESOLVE", constants.TargetDelay, fmt.Sprint(id),
		utils.RenderLog(constants.LogTemplates["DELAY_RESOLVE"], map[string]interface{}{
			"DelayID": id, "ResolvedAt": now.Format(time.RFC3339), "Minutes": event.Minutes,
		}))
	return &view, nil
}

func (s *DelayEventService) currentSummary(turnaroundID uint) (utils.OpenDelaySummary, error) {
	delayMap, err := s.delayRepo.MapByTurnaround([]uint{turnaroundID})
	if err != nil {
		return utils.OpenDelaySummary{}, err
	}
	return utils.SummarizeOpenDelays(delayMap[turnaroundID]), nil
}
