import { StatusBadge } from "./StatusBadge";
import type { GroundTask } from "../../types/GroundTask";
import { formatClock, formatOverdue } from "../../utils/formatters";

/**
 * 超时任务卡片 —— 看板与任务列表共用的判定结果展示。
 * 数据中的 overdue / overdue_minutes / effective_deadline 由 utils/delayPolicy（前端）
 * 或 ViewService（后端）统一算出，本组件不自行判超时。
 */
export function OverdueTaskCard({ task }: { task: GroundTask }) {
  return (
    <article className="overdue-card">
      <header>
        <strong>#{task.id} · {task.flight_no} · {task.task_type_text}</strong>
        <StatusBadge value="OVERDUE" tone="danger" text={formatOverdue(task.overdue_minutes, task.overdue)} />
      </header>
      <div className="overdue-grid">
        <span>原计划截止</span><time>{formatClock(task.deadline)}</time>
        <span>顺延 {task.open_delay_minutes} 分钟</span>
        <time>{formatClock(task.effective_deadline)}</time>
        <span>比较基准</span>
        <span>{task.compare_basis === "ACTUAL_FINISH" ? `实际完成 ${formatClock(task.actual_finish)}` : "当前时间"}</span>
      </div>
      {task.blocker_note ? <p className="blocker">阻塞：{task.blocker_note}</p> : null}
    </article>
  );
}
