import { useEffect, useMemo } from "react";
import { useGroundTaskStore } from "../stores/GroundTaskStore";
import { useDelayEventStore } from "../stores/DelayEventStore";
import { useFlightTurnaroundStore } from "../stores/FlightTurnaroundStore";
import { useTaskOvertime } from "../hooks/useTaskOvertime";
import { formatDateOrDash, formatMinutes } from "../utils/formatters";
import { STATUS_TEXT } from "../constants/statusText";
import { StatCard } from "../components/common/StatCard";
import { StatusBadge } from "../components/common/StatusBadge";
import { DelayTag } from "../components/common/DelayTag";
import { OvertimeRuleNote } from "../components/common/OvertimeRuleNote";
import { EmptyState } from "../components/common/EmptyState";

const taskTypeText = (value: string) =>
  (STATUS_TEXT.GroundTaskType as Record<string, string>)[value] ?? value;

export function DashboardPage() {
  const loadTasks = useGroundTaskStore((state) => state.load);
  const loadDelays = useDelayEventStore((state) => state.load);
  const turnarounds = useFlightTurnaroundStore((state) => state.rows);
  const loadTurnarounds = useFlightTurnaroundStore((state) => state.load);

  useEffect(() => {
    void loadTasks();
    void loadDelays();
    void loadTurnarounds();
  }, [loadTasks, loadDelays, loadTurnarounds]);

  // 超时卡片与任务列表共用同一套口径（useTaskOvertime）
  const { rows, overtimeRows, overtimeCount, openDelayCount, openDelayMinutesTotal, evaluatedAt } =
    useTaskOvertime();

  const flightNo = useMemo(() => {
    const map = new Map(turnarounds.map((item) => [item.id, item.flight_no]));
    return (id: number) => map.get(id) ?? `#${id}`;
  }, [turnarounds]);

  const signedCount = rows.filter(({ task }) => task.actual_finish).length;

  return (
    <main className="page">
      <section className="page-head">
        <div>
          <p className="eyebrow">ground-turn</p>
          <h1>过站运行看板</h1>
        </div>
        <StatusBadge value={`评估 ${formatDateOrDash(evaluatedAt.toISOString())}`} />
      </section>

      <section className="metrics">
        <StatCard label="超时任务（按顺延后生效截止）" value={overtimeCount} />
        <StatCard label="未关闭延误事件" value={openDelayCount} />
        <StatCard label="未关闭延误顺延分钟" value={formatMinutes(openDelayMinutesTotal)} />
        <StatCard label="已签收 / 任务总数" value={`${signedCount} / ${rows.length}`} />
      </section>

      <OvertimeRuleNote />

      <section className="panel wide">
        <h2>超时任务明细</h2>
        {overtimeRows.length === 0 ? (
          <EmptyState title="当前没有超时任务" />
        ) : (
          <table className="data-table">
            <thead>
              <tr>
                <th>航班</th>
                <th>任务类型</th>
                <th>原计划截止</th>
                <th>延误顺延</th>
                <th>当前生效截止</th>
                <th>签收时间</th>
                <th>超时</th>
              </tr>
            </thead>
            <tbody>
              {overtimeRows.map(({ task, delayMinutes, effectiveDeadline, overtimeMinutes }) => (
                <tr key={task.id} className="overtime">
                  <td>{flightNo(task.turnaround_id)}</td>
                  <td>{taskTypeText(task.task_type)}</td>
                  <td>{formatDateOrDash(task.deadline)}</td>
                  <td>{delayMinutes > 0 ? <DelayTag minutes={delayMinutes} /> : "—"}</td>
                  <td>{formatDateOrDash(effectiveDeadline)}</td>
                  <td>{formatDateOrDash(task.actual_finish)}</td>
                  <td>超时 {formatMinutes(overtimeMinutes)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </section>
    </main>
  );
}
