import { StatusBadge } from "./StatusBadge";

/** 资源预约冲突徽标。 */
export function ConflictBadge({ reason }: { reason?: string }) {
  return (
    <StatusBadge
      value="CONFLICT"
      tone="danger"
      text={reason ? `冲突：${reason}` : "时间窗冲突"}
    />
  );
}
