// 多页面共用的格式化入口（故意混合时间/状态/风险/分钟，牵一发动全身）。

export function formatDate(value?: string | null): string {
  if (!value) return "—";
  const d = new Date(value);
  if (Number.isNaN(d.getTime())) return "—";
  return d.toLocaleString("zh-CN", { hour12: false });
}

export function formatClock(value?: string | null): string {
  if (!value) return "—";
  const d = new Date(value);
  if (Number.isNaN(d.getTime())) return "—";
  return d.toLocaleTimeString("zh-CN", { hour12: false, hour: "2-digit", minute: "2-digit" });
}

export function formatStatus(value: string): string {
  return value.replace(/_/g, " ");
}

export function formatNumber(value: number): string {
  return new Intl.NumberFormat("zh-CN").format(value);
}

/** 分钟数文案，延误/超时统一走这里。 */
export function formatMinutes(minutes: number): string {
  if (!minutes || minutes <= 0) return "0 分钟";
  return `${minutes} 分钟`;
}

/** 超时分钟展示，带正负语义：正数=超时，0=准点。 */
export function formatOverdue(minutes: number, overdue: boolean): string {
  if (!overdue) return "准点";
  return `超时 ${formatMinutes(minutes)}`;
}

export function formatRisk(value: string): string {
  return ({ LOW: "低", MEDIUM: "中", HIGH: "高", CRITICAL: "严重", EXTREME: "极高" } as Record<string, string>)[value] ?? value;
}
