import { StatusBadge } from "./StatusBadge";
import { formatMinutes } from "../../utils/formatters";

// 延误标签：只用于“延误分钟”的统一展示；closed 置灰，避免把已关闭延误误读为仍在顺延。
export function DelayTag({
  minutes,
  closed,
  label
}: {
  minutes: number;
  closed?: boolean;
  label?: string;
}) {
  return (
    <StatusBadge
      value={closed ? "DELAY_CLOSED" : "DELAY_OPEN"}
      tone={closed ? "muted" : "warn"}
      text={`${label ?? "延误"} ${formatMinutes(minutes)}${closed ? "（已关闭·不顺延）" : ""}`}
    />
  );
}
