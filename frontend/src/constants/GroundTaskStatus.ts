export const GroundTaskStatus = ["PENDING", "SIGNED", "DONE", "BLOCKED"] as const;
export type GroundTaskStatus = (typeof GroundTaskStatus)[number];
export const GroundTaskStatusText: Record<GroundTaskStatus, string> = {
  PENDING: "待签收",
  SIGNED: "已签收",
  DONE: "已完成",
  BLOCKED: "阻塞"
};
