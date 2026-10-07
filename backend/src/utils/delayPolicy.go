package utils

import (
	"fmt"
	"strings"
	"time"

	"groundTurn/src/models"
)

// 本文件是「延误顺延 -> 超时判定」的唯一口径来源，
// 后端 service、看板卡片、任务列表都调用这里；前端 utils/delayPolicy.ts 保持同构实现。

// OpenDelaySummary 汇总某航班上尚未关闭的延误事件。
// 关闭过（resolved_at 非空）的延误不参与顺延。
type OpenDelaySummary struct {
	OpenMinutes int
	OpenCount   int
}

// SummarizeOpenDelays 只统计 ResolvedAt == nil 的事件。
func SummarizeOpenDelays(events []models.DelayEvent) OpenDelaySummary {
	s := OpenDelaySummary{}
	for i := range events {
		if events[i].ResolvedAt == nil {
			s.OpenMinutes += events[i].Minutes
			s.OpenCount++
		}
	}
	return s
}

// EffectiveDeadline 保留原计划 deadline，另算当前生效截止：
//
//	effective_deadline = deadline + Σ 未关闭延误分钟
//
// 登记延误时绝不改写任务 deadline 列，只在读取时重算。
func EffectiveDeadline(deadline time.Time, summary OpenDelaySummary) time.Time {
	return deadline.Add(time.Duration(summary.OpenMinutes) * time.Minute)
}

// OverdueVerdict 超时判定结果。
type OverdueVerdict struct {
	Overdue        bool
	OverdueMinutes int
	// CompareBasis: NOW=未完成，按当前时间；ACTUAL_FINISH=已签收/完成，按实际完成时间。
	CompareBasis string
}

const (
	BasisNow          = "NOW"
	BasisActualFinish = "ACTUAL_FINISH"
)

// EvaluateOverdue 任务超时一律按顺延后的生效截止 effectiveDeadline 判定，
// 不再使用原计划 deadline。
//   - 已完成（actual_finish 有值，含签收后完成）：actual_finish > effective_deadline 即超时
//   - 未完成：now > effective_deadline 即超时
func EvaluateOverdue(status string, effectiveDeadline time.Time, actualFinish *time.Time, now time.Time) OverdueVerdict {
	if actualFinish != nil {
		return OverdueVerdict{
			Overdue:        actualFinish.After(effectiveDeadline),
			OverdueMinutes: overdueMinutes(*actualFinish, effectiveDeadline),
			CompareBasis:   BasisActualFinish,
		}
	}
	_ = status
	return OverdueVerdict{
		Overdue:        now.After(effectiveDeadline),
		OverdueMinutes: overdueMinutes(now, effectiveDeadline),
		CompareBasis:   BasisNow,
	}
}

func overdueMinutes(t time.Time, deadline time.Time) int {
	if !t.After(deadline) {
		return 0
	}
	d := t.Sub(deadline)
	return int(d.Round(time.Minute) / time.Minute)
}

// RenderLog 渲染 constants.LogTemplates 中的 {{.Field}} 模板。
func RenderLog(template string, fields map[string]interface{}) string {
	out := template
	for k, v := range fields {
		out = strings.ReplaceAll(out, "{{."+k+"}}", fmt.Sprintf("%v", v))
	}
	return out
}
