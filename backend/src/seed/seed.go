package seed

import (
	"time"

	"groundTurn/src/config"
	"groundTurn/src/constants"
	"groundTurn/src/models"
	"groundTurn/src/utils"

	"gorm.io/gorm"
)

// Run 写入本地演示数据（幂等：已有用户则跳过）。
// 所有时间相对当前时刻生成，保证看板超时卡片在任何启动时间都有演示效果。
func Run(db *gorm.DB, cfg config.Config) error {
	var userCount int64
	if err := db.Model(&models.User{}).Count(&userCount).Error; err != nil {
		return err
	}
	if userCount > 0 {
		return nil
	}

	now := time.Now().UTC()
	at := func(min int) time.Time { return now.Add(time.Duration(min) * time.Minute) }
	ptr := func(t time.Time) *time.Time { return &t }

	// ---- 账号（口令统一 demo123）----
	hash := func(pw string) string { return utils.HashPassword(pw, cfg.JWTSecret) }
	users := []models.User{
		{Username: "dispatcher", PasswordHash: hash("demo123"), Role: constants.RoleDispatcher, TeamCode: "OPS", DisplayName: "调度-张运"},
		{Username: "team", PasswordHash: hash("demo123"), Role: constants.RoleTeam, TeamCode: "RED", DisplayName: "红班-李勤"},
		{Username: "resource", PasswordHash: hash("demo123"), Role: constants.RoleResourceManager, TeamCode: "EQP", DisplayName: "资源-王管"},
		{Username: "supervisor", PasswordHash: hash("demo123"), Role: constants.RoleSupervisor, TeamCode: "OPS", DisplayName: "督导-赵督"},
	}
	if err := db.Create(&users).Error; err != nil {
		return err
	}

	// ---- 航班 ----
	flights := []models.FlightTurnaround{
		{FlightNo: "CA1234", AircraftReg: "B-5861", StandNo: "201", ArrivalTime: at(-90), DepartureTime: at(30), TurnaroundStatus: constants.TurnaroundInService, DelayReason: "天气+餐食延误"},
		{FlightNo: "CA5678", AircraftReg: "B-6207", StandNo: "118", ArrivalTime: at(-20), DepartureTime: at(40), TurnaroundStatus: constants.TurnaroundOnStand},
		{FlightNo: "CA9012", AircraftReg: "B-3099", StandNo: "305", ArrivalTime: at(-160), DepartureTime: at(-60), TurnaroundStatus: constants.TurnaroundDeparted},
	}
	if err := db.Create(&flights).Error; err != nil {
		return err
	}

	// ---- 延误事件 ----
	// 航班1：未关闭 20+10=30 分钟（顺延）；已关闭 40 分钟（不顺延）。
	// 航班2：未关闭 25 分钟（顺延）。
	// 航班3：仅一笔已关闭 15 分钟（不顺延）。
	delays := []models.DelayEvent{
		{TurnaroundID: flights[0].ID, DelayType: "WEATHER", Minutes: 20, RootCause: "本场雷雨", ResponsibilityTeam: "OPS", CreatedAt: at(-70), ResolvedAt: nil},
		{TurnaroundID: flights[0].ID, DelayType: "CATERING", Minutes: 10, RootCause: "餐车晚到", ResponsibilityTeam: "CAT", CreatedAt: at(-50), ResolvedAt: nil},
		{TurnaroundID: flights[0].ID, DelayType: "BAGGAGE", Minutes: 40, RootCause: "传送带故障（已修复）", ResponsibilityTeam: "BAG", CreatedAt: at(-80), ResolvedAt: ptr(at(-45))},
		{TurnaroundID: flights[1].ID, DelayType: "ATC_FLOW", Minutes: 25, RootCause: "流量控制", ResponsibilityTeam: "ATC", CreatedAt: at(-15), ResolvedAt: nil},
		{TurnaroundID: flights[2].ID, DelayType: "REFUEL", Minutes: 15, RootCause: "加油单核对（已闭环）", ResponsibilityTeam: "FUL", CreatedAt: at(-140), ResolvedAt: ptr(at(-120))},
	}
	if err := db.Create(&delays).Error; err != nil {
		return err
	}

	// ---- 任务（deadline 一律是原计划截止，绝不写入顺延值）----
	tasks := []models.GroundTask{
		// 航班1 生效截止 = 原计划 + 30min
		{TurnaroundID: flights[0].ID, TaskType: "CATERING", TeamID: "CAT", PlannedStart: ptr(at(-60)), Deadline: at(-15), Status: constants.TaskStatusPending},
		{TurnaroundID: flights[0].ID, TaskType: "BAGGAGE", TeamID: "BAG", PlannedStart: ptr(at(-80)), Deadline: at(-45), AcceptedAt: ptr(at(-75)), ActualFinish: ptr(at(-40)), Status: constants.TaskStatusCompleted},
		{TurnaroundID: flights[0].ID, TaskType: "CLEANING", TeamID: "CLN", PlannedStart: ptr(at(-85)), Deadline: at(-50), AcceptedAt: ptr(at(-80)), ActualFinish: ptr(at(-10)), Status: constants.TaskStatusCompleted},
		{TurnaroundID: flights[0].ID, TaskType: "REFUEL", TeamID: "FUL", PlannedStart: ptr(at(-30)), Deadline: at(-5), Status: constants.TaskStatusPending},
		{TurnaroundID: flights[0].ID, TaskType: "PUSHBACK", TeamID: "PUB", PlannedStart: ptr(at(-45)), Deadline: at(-40), Status: constants.TaskStatusInProgress, AcceptedAt: ptr(at(-42))},

		// 航班2 生效截止 = 原计划 + 25min
		{TurnaroundID: flights[1].ID, TaskType: "CATERING", TeamID: "CAT", PlannedStart: ptr(at(-10)), Deadline: at(10), Status: constants.TaskStatusPending},
		{TurnaroundID: flights[1].ID, TaskType: "WATER_SERVICE", TeamID: "WTR", PlannedStart: ptr(at(-25)), Deadline: at(-30), Status: constants.TaskStatusInProgress, AcceptedAt: ptr(at(-25))},

		// 航班3 已离港，无未关闭延误，生效截止=原计划
		{TurnaroundID: flights[2].ID, TaskType: "BAGGAGE", TeamID: "BAG", PlannedStart: ptr(at(-150)), Deadline: at(-120), AcceptedAt: ptr(at(-148)), ActualFinish: ptr(at(-130)), Status: constants.TaskStatusCompleted},
		{TurnaroundID: flights[2].ID, TaskType: "REFUEL", TeamID: "FUL", PlannedStart: ptr(at(-130)), Deadline: at(-100), AcceptedAt: ptr(at(-128)), ActualFinish: ptr(at(-90)), Status: constants.TaskStatusCompleted},
	}
	if err := db.Create(&tasks).Error; err != nil {
		return err
	}

	// ---- 资源 ----
	resources := []models.GroundResource{
		{ResourceCode: "GPU-01", ResourceType: "POWER", Location: "201", AvailabilityStatus: constants.ResourceBooked, OwnerTeam: "EQP"},
		{ResourceCode: "CAT-07", ResourceType: "CATERING", Location: "厨房-2", AvailabilityStatus: constants.ResourceAvailable, OwnerTeam: "CAT"},
		{ResourceCode: "FUEL-03", ResourceType: "REFUEL", Location: "油站-A", AvailabilityStatus: constants.ResourceMaintenance, MaintenanceDueAt: ptr(at(240)), OwnerTeam: "FUL"},
		{ResourceCode: "BLT-12", ResourceType: "BAGGAGE", Location: "分拣-1", AvailabilityStatus: constants.ResourceOffline, OwnerTeam: "BAG"},
	}
	if err := db.Create(&resources).Error; err != nil {
		return err
	}

	// ---- 预约（含一对时间冲突，供资源页 ConflictBadge）----
	bookings := []models.ResourceBooking{
		{ResourceID: resources[0].ID, TurnaroundID: flights[0].ID, TaskID: &tasks[0].ID, StartTime: at(-60), EndTime: at(0), BookingStatus: constants.BookingConfirmed},
		{ResourceID: resources[0].ID, TurnaroundID: flights[1].ID, StartTime: at(-30), EndTime: at(20), BookingStatus: constants.BookingConflict, ConflictReason: "overlaps booking on GPU-01"},
		{ResourceID: resources[1].ID, TurnaroundID: flights[1].ID, TaskID: &tasks[5].ID, StartTime: at(-10), EndTime: at(20), BookingStatus: constants.BookingHeld},
		{ResourceID: resources[2].ID, TurnaroundID: flights[0].ID, StartTime: at(-40), EndTime: at(-10), BookingStatus: constants.BookingReleased},
	}
	if err := db.Create(&bookings).Error; err != nil {
		return err
	}

	return nil
}
