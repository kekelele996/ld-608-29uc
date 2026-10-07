import { useEffect, useMemo } from "react";
import { useGroundTaskStore } from "../stores/GroundTaskStore";
import { useDelayEventStore } from "../stores/DelayEventStore";
import { useFlightTurnaroundStore } from "../stores/FlightTurnaroundStore";
import { useTaskOvertime } from "../hooks/useTaskOvertime";
import { hasActualFinish } from "../utils/overtime";
import { formatDateOrDash, formatMinutes } from "../utils/formatters";
import { STATUS_TEXT } from "../constants/statusText";
import { StatusBadge } from "../components/common/StatusBadge";
import { DelayTag } from "../components/common/DelayTag";
import { OvertimeRuleNote } from "../components/common/OvertimeRuleNote";
import { EmptyState } from "../components/common/EmptyState";

const taskTypeText = (value: string) =>
  (STATUS_TEXT.GroundTaskType as Record<string, string>)[value] ?? value;
const taskStatusText = (value: string) =>
  (STATUS_TEXT.GroundTaskStatus as Record<string, string>)[value] ?? value;

export function TasksPage() {
  const tasksLoading = useGroundTaskStore((state) => state.loading);
  const loadTasks = useGroundTaskStore((state) => state.load);
  const signOff = useGroundTaskStore((state) => state.signOff);
  const loadDelays = useDelayEventStore((state) => state.load);
  const turnarounds = useFlightTurnaroundStore((state) => state.rows);
  const loadTurnarounds = useFlightTurnaroundStore((state) => state.load);

  useEffect(() => {
    void loadTasks();
    void loadDelays();
    void loadTurnarounds();
  }, [loadTasks, loadDelays, loadTurnarounds]);

  // 与看板超时卡片共用同一套口径
  const { rows, overtimeCount, evaluatedAt } = useTaskOvertime();

  const flightNo = useMemo(() => {
    const map = new Map(turnarounds.map((item) => [item.id, item.flight_no]));
    return (id: number) => map.get(id) ?? `#${id}`;
  }, [turnarounds]);

  return (
    <main className="page">
      <section className="page-head">
        <div>
          <p className="eyebrow">ground-turn</p>
          <h1>地勤任务</h1>
        </div>
        <StatusBadge value={`超时 ${overtimeCount} 条`} />
      </section>

      <OvertimeRuleNote />

      <section className="panel wide">
        <h2>任务列表（评估时刻 {formatDateOrDash(evaluatedAt.toISOString())}）</h2>
        {rows.length === 0 && !tasksLoading ? (
          <EmptyState title="暂无任务" />
        ) : (
          <table className="data-table">
            <thead>
              <tr>
                <th>航班</th>
                <th>任务类型</th>
                <th>班组</th>
                <th>原计划截止</th>
                <th>延误顺延</th>
                <th>当前生效截止</th>
                <th>状态</th>
                <th>签收时间</th>
                <th>超时判定</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              {rows.map(({ task, delayMinutes, effectiveDeadline, overtimeMinutes, isOvertime }) => (
                <tr key={task.id} className={isOvertime ? "overtime" : ""}>
                  <td>{flightNo(task.turnaround_id)}</td>
                  <td>{taskTypeText(task.task_type)}</td>
                  <td>班组 {task.team_id}</td>
                  <td>{formatDateOrDash(task.deadline)}</td>
                  <td>{delayMinutes > 0 ? <DelayTag minutes={delayMinutes} /> : "—"}</td>
                  <td>{formatDateOrDash(effectiveDeadline)}</td>
                  <td>
                    <StatusBadge value={task.status} /> {taskStatusText(task.status)}
                  </td>
                  <td>{formatDateOrDash(task.actual_finish)}</td>
                  <td>{isOvertime ? `超时 ${formatMinutes(overtimeMinutes)}` : "未超时"}</td>
                  <td>
                    {hasActualFinish(task) ? (
                      <span className="muted">已签收</span>
                    ) : (
                      <button className="btn primary" onClick={() => void signOff(task.id)}>
                        签收
                      </button>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </section>
    </main>
  );
}
