package constructors

import (
	"testing"
	"time"

	"groundTurn/src/constants"
	"groundTurn/src/models"
	"groundTurn/src/types"
	"groundTurn/src/utils"
)

// 以固定 now 构造：航班有未关闭延误 20+10=30，已关闭延误 40（不顺延）。
func TestBuildGroundTaskViewDeadlinePolicy(t *testing.T) {
	now := time.Unix(0, 0)
	deadline := now.Add(-15 * time.Minute) // 原计划截止在 15 分钟前
	summary := utils.OpenDelaySummary{OpenMinutes: 30, OpenCount: 2}

	// 1) 未完成任务：生效截止 = -15+30 = +15 分钟后，未超时；按当前时间
	pending := models.GroundTask{ID: 1, TurnaroundID: 1, TaskType: "CATERING", Deadline: deadline, Status: constants.TaskStatusPending}
	v1 := BuildGroundTaskView(pending, "CA1", summary, now)
	if v1.Overdue {
		t.Fatalf("pending should be saved by +30m extension, got overdue=%+v", v1)
	}
	if v1.CompareBasis != utils.BasisNow {
		t.Fatalf("pending compare basis = %s, want NOW", v1.CompareBasis)
	}
	if !v1.EffectiveDeadline.Equal(deadline.Add(30 * time.Minute)) {
		t.Fatalf("effective = %v, want deadline+30m", v1.EffectiveDeadline)
	}
	if !v1.Deadline.Equal(deadline) {
		t.Fatal("deadline column must remain the original plan")
	}

	// 2) 已完成，actual_finish 在 +20（晚于生效截止 +15）：超时 5 分钟，按实际完成时间
	finish := now.Add(20 * time.Minute)
	done := models.GroundTask{ID: 2, TurnaroundID: 1, TaskType: "BAGGAGE", Deadline: deadline,
		AcceptedAt: &now, ActualFinish: &finish, Status: constants.TaskStatusCompleted}
	v2 := BuildGroundTaskView(done, "CA1", summary, now)
	if !v2.Overdue || v2.OverdueMinutes != 5 || v2.CompareBasis != utils.BasisActualFinish {
		t.Fatalf("completed late = %+v", v2)
	}

	// 3) 若延误全部关闭（顺延归零，生效截止=原计划 -15），actual(+20) 超时 35 分钟
	closedSummary := utils.OpenDelaySummary{OpenMinutes: 0, OpenCount: 0}
	v3 := BuildGroundTaskView(done, "CA1", closedSummary, now)
	if !v3.Overdue || v3.OverdueMinutes != 35 {
		t.Fatalf("after delays closed, completed task overdue minutes = %+v, want 35", v3)
	}
}

// 航班汇总：超时条数由任务视图逐个数出，必须与顺延分钟口径一致。
func TestFlightViewOverdueCountMatchesTasks(t *testing.T) {
	now := time.Unix(0, 0)
	deadline := now.Add(-15 * time.Minute)
	summary := utils.OpenDelaySummary{OpenMinutes: 30, OpenCount: 2}

	taskViews := []types.GroundTaskView{
		BuildGroundTaskView(models.GroundTask{ID: 1, TaskType: "CATERING", Deadline: deadline, Status: constants.TaskStatusPending}, "CA1", summary, now),
		BuildGroundTaskView(models.GroundTask{ID: 2, TaskType: "CLEANING", Deadline: deadline, Status: constants.TaskStatusInProgress}, "CA1", summary, now),
	}
	finish := now.Add(40 * time.Minute)
	taskViews = append(taskViews, BuildGroundTaskView(models.GroundTask{
		ID: 3, TaskType: "BAGGAGE", Deadline: deadline, ActualFinish: &finish, Status: constants.TaskStatusCompleted,
	}, "CA1", summary, now))

	overdue := 0
	for _, tv := range taskViews {
		if tv.Overdue {
			overdue++
		}
	}
	if overdue != 1 {
		t.Fatalf("overdue count = %d, want 1 (reconciles with delay minutes)", overdue)
	}

	flight := models.FlightTurnaround{ID: 1, FlightNo: "CA1", TurnaroundStatus: constants.TurnaroundInService}
	fv := BuildFlightTurnaroundView(flight, taskViews, summary)
	if fv.TaskOverdue != 1 || fv.TaskTotal != 3 || fv.TaskCompleted != 1 || fv.OpenDelayMinutes != 30 {
		t.Fatalf("flight view aggregates mismatch: %+v", fv)
	}
}
