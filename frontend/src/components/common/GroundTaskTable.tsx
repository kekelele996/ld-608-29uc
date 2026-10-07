import type { GroundTask } from "../../types/GroundTask";
import { StatusBadge } from "./StatusBadge";
import { DelayTag } from "./DelayTag";
import { formatClock, formatDate, formatOverdue } from "../../utils/formatters";

interface Props {
  tasks: GroundTask[];
  canSign?: boolean;
  onAccept?: (id: number) => void;
  onComplete?: (id: number) => void;
}

function statusTone(task: GroundTask): "danger" | "warn" | "ok" | "muted" {
  if (task.status === "BLOCKED") return "warn";
  if (task.status === "COMPLETED") return "ok";
  if (task.status === "PENDING") return "muted";
  return "ok";
}

/**
 * 任务列表（任务页主表，看板“超时任务”之外的明细也复用其行内字段）。
 * 截止列同时展示原计划与生效截止，超时光标高亮的是生效截止。
 */
export function GroundTaskTable({ tasks, canSign, onAccept, onComplete }: Props) {
  return (
    <div className="table-wrap">
      <table className="data-table">
        <thead>
          <tr>
            <th>#</th><th>航班</th><th>任务</th><th>班组</th>
            <th>原计划截止</th><th>未关闭延误</th><th>当前生效截止</th>
            <th>签收时间</th><th>实际完成</th><th>状态</th><th>超时（按生效截止）</th>
            {(canSign && (onAccept || onComplete)) ? <th>操作</th> : null}
          </tr>
        </thead>
        <tbody>
          {tasks.map((t) => (
            <tr key={t.id} className={t.overdue ? "row-overdue" : ""}>
              <td>{t.id}</td>
              <td>{t.flight_no}</td>
              <td>{t.task_type_text}</td>
              <td>{t.team_id || "—"}</td>
              <td>{formatClock(t.deadline)}</td>
              <td>{t.open_delay_minutes > 0 ? <DelayTag minutes={t.open_delay_minutes} label="" /> : <span>0 分钟</span>}</td>
              <td><strong>{formatClock(t.effective_deadline)}</strong></td>
              <td>{t.accepted_at ? formatDate(t.accepted_at) : <em className="dim">未签收</em>}</td>
              <td>{t.actual_finish ? formatDate(t.actual_finish) : <em className="dim">—</em>}</td>
              <td><StatusBadge value={t.status} text={t.status_text} tone={statusTone(t)} /></td>
              <td>
                {t.overdue
                  ? <StatusBadge value="OVERDUE" tone="danger" text={formatOverdue(t.overdue_minutes, true)} />
                  : <StatusBadge value="ON_TIME" tone="ok" text="准点" />}
                <div className="dim basis">按{t.compare_basis === "ACTUAL_FINISH" ? "实际完成" : "当前时间"}比较</div>
              </td>
              {canSign && (onAccept || onComplete) ? (
                <td className="actions">
                  {onAccept && !t.accepted_at && t.status !== "COMPLETED"
                    ? <button className="btn small" onClick={() => onAccept(t.id)}>签收</button> : null}
                  {onComplete && t.status !== "COMPLETED"
                    ? <button className="btn small primary" onClick={() => onComplete(t.id)}>完成</button> : null}
                </td>
              ) : null}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
