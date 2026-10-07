export const GroundTaskType = ["CLEANING", "CATERING", "BAGGAGE", "REFUEL", "WATER_SERVICE", "PUSHBACK"] as const;
export type GroundTaskType = (typeof GroundTaskType)[number];
export const GroundTaskTypeText: Record<GroundTaskType, string> = {
  CLEANING: "客舱清洁",
  CATERING: "餐食保障",
  BAGGAGE: "行李装卸",
  REFUEL: "加油",
  WATER_SERVICE: "清水污水",
  PUSHBACK: "推机"
};
