export const ResourceStatus = ["AVAILABLE", "BOOKED", "MAINTENANCE", "OFFLINE"] as const;
export type ResourceStatus = (typeof ResourceStatus)[number];
export const ResourceStatusText: Record<ResourceStatus, string> = {
  AVAILABLE: "可用",
  BOOKED: "已预约",
  MAINTENANCE: "检修中",
  OFFLINE: "下线"
};

export const BookingStatus = ["HELD", "CONFIRMED", "RELEASED", "CONFLICT"] as const;
export type BookingStatus = (typeof BookingStatus)[number];
export const BookingStatusText: Record<BookingStatus, string> = {
  HELD: "占用",
  CONFIRMED: "已确认",
  RELEASED: "已释放",
  CONFLICT: "冲突"
};
