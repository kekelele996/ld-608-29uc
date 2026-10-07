// 全局共用状态徽标。overdue 时变红，便于看板卡片与任务列表共用同一视觉口径。
export function StatusBadge({ value, text, tone }: { value: string; text?: string; tone?: "danger" | "warn" | "ok" | "muted" }) {
  const cls = tone ? `badge tone-${tone}` : "badge " + String(value).toLowerCase().replace(/_/g, "-");
  return <span className={cls}>{text ?? String(value).replace(/_/g, " ")}</span>;
}
