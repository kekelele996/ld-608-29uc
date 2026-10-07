export const formatDate = (value: string) => new Date(value).toLocaleString("zh-CN", { hour12: false });
export const formatDateOrDash = (value: string) => {
  const time = Date.parse(value);
  return Number.isNaN(time) ? "—" : new Date(time).toLocaleString("zh-CN", { hour12: false });
};
export const formatStatus = (value: string) => value.replace(/_/g, " ");
export const formatNumber = (value: number) => new Intl.NumberFormat("zh-CN").format(value);
export const formatMinutes = (value: number) => `${formatNumber(value)} 分钟`;
export const formatRisk = (value: string) => ({ LOW: "低", MEDIUM: "中", HIGH: "高", CRITICAL: "严重", EXTREME: "极高" }[value] ?? value);
