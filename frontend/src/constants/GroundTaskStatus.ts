export const GroundTaskStatus = ["PENDING", "ACCEPTED", "IN_PROGRESS", "BLOCKED", "COMPLETED"] as const;
export type GroundTaskStatus = (typeof GroundTaskStatus)[number];

export const GroundTaskStatusText: Record<GroundTaskStatus, string> = {
  PENDING: "待签收",
  ACCEPTED: "已签收",
  IN_PROGRESS: "进行中",
  BLOCKED: "阻塞",
  COMPLETED: "已完成"
};

// 角色（与后端 RBAC 对齐）
export const Roles = ["DISPATCHER", "TEAM", "RESOURCE_MANAGER", "SUPERVISOR"] as const;
export type Role = (typeof Roles)[number];

export const RoleText: Record<Role, string> = {
  DISPATCHER: "地勤调度",
  TEAM: "班组",
  RESOURCE_MANAGER: "资源管理员",
  SUPERVISOR: "运行督导"
};

// 延误类型
export const DelayType = ["CLEANING", "CATERING", "BAGGAGE", "REFUEL", "WATER_SERVICE", "PUSHBACK", "WEATHER", "ATC_FLOW"] as const;
export type DelayType = (typeof DelayType)[number];

export const DelayTypeText: Record<DelayType, string> = {
  CLEANING: "客舱清洁",
  CATERING: "餐食保障",
  BAGGAGE: "行李装卸",
  REFUEL: "加油",
  WATER_SERVICE: "清水污水",
  PUSHBACK: "推机",
  WEATHER: "天气",
  ATC_FLOW: "流量控制"
};
