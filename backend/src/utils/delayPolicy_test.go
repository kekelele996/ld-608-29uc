package utils

import (
	"testing"
	"time"

	"groundTurn/src/models"
)

func tm(min int) time.Time { return time.Unix(0, 0).Add(time.Duration(min) * time.Minute) }

func TestSummarizeOpenDelays(t *testing.T) {
	resolved := tm(100)
	events := []models.DelayEvent{
		{Minutes: 20, ResolvedAt: nil},       // 未关闭
		{Minutes: 10, ResolvedAt: nil},       // 未关闭
		{Minutes: 40, ResolvedAt: &resolved}, // 已关闭：不计
	}
	got := SummarizeOpenDelays(events)
	if got.OpenMinutes != 30 || got.OpenCount != 2 {
		t.Fatalf("open summary = %+v, want 30min/2", got)
	}
}

func TestEffectiveDeadlineKeepsOriginalAndExtends(t *testing.T) {
	base := tm(0)
	eff := EffectiveDeadline(base, OpenDelaySummary{OpenMinutes: 30})
	if !eff.Equal(tm(30)) {
		t.Fatalf("effective deadline = %v, want +30m", eff)
	}
	// 原 deadline 不变（函数纯计算）
	if !base.Equal(tm(0)) {
		t.Fatal("original deadline must be preserved")
	}
}

func TestEvaluateOverdue(t *testing.T) {
	eff := tm(30)
	now := tm(35)

	// 未完成，now(35) > effective(30)：超时 5 分钟，按当前时间
	v := EvaluateOverdue("PENDING", eff, nil, now)
	if !v.Overdue || v.OverdueMinutes != 5 || v.CompareBasis != BasisNow {
		t.Fatalf("pending overdue = %+v", v)
	}

	// 未完成但 now 未超过：不超时
	v2 := EvaluateOverdue("PENDING", eff, nil, tm(20))
	if v2.Overdue {
		t.Fatalf("should be on time, got %+v", v2)
	}

	// 已完成 actual(40) > effective(30)：超时 10，按实际完成时间
	finish := tm(40)
	v3 := EvaluateOverdue("COMPLETED", eff, &finish, now)
	if !v3.Overdue || v3.OverdueMinutes != 10 || v3.CompareBasis != BasisActualFinish {
		t.Fatalf("completed overdue = %+v", v3)
	}

	// 已完成但在生效截止内（如延误已关闭、按短截止也完成得早）：不超时
	early := tm(25)
	v4 := EvaluateOverdue("COMPLETED", eff, &early, now)
	if v4.Overdue || v4.CompareBasis != BasisActualFinish {
		t.Fatalf("completed on time = %+v", v4)
	}
}

// 关键回归：关闭延误后，顺延分钟回落，原本被顺延“救回”的任务重新超时，
// 且超时条数随顺延分钟变化一致。
func TestClosingDelayChangesOverdueCount(t *testing.T) {
	deadline := tm(0)
	now := tm(20) // 原 deadline 已过 20 分钟

	open := OpenDelaySummary{OpenMinutes: 30}
	if EvaluateOverdue("PENDING", EffectiveDeadline(deadline, open), nil, now).Overdue {
		t.Fatal("with +30m extension task should be within effective deadline")
	}
	closed := OpenDelaySummary{OpenMinutes: 0}
	if !EvaluateOverdue("PENDING", EffectiveDeadline(deadline, closed), nil, now).Overdue {
		t.Fatal("after delays closed (+0m) task must be overdue against original deadline")
	}
}
