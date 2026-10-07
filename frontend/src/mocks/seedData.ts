import type { DelayEvent } from "../types/DelayEvent";
import type { FlightTurnaround } from "../types/FlightTurnaround";
import type { GroundResource } from "../types/GroundResource";
import type { GroundTask } from "../types/GroundTask";
import type { ResourceBooking } from "../types/ResourceBooking";
import { decorateTasks } from "../utils/delayPolicy";
import { GroundTaskTypeText } from "../constants/GroundTaskType";
import { GroundTaskStatusText } from "../constants/GroundTaskStatus";
import { TurnaroundStatusText } from "../constants/TurnaroundStatus";

/**
 * 本地种子数据（禁止第三方 API）。后端不可达时作为离线兜底。
 * 时间相对当前时刻，保证超时演示稳定；任务的 effective/overdue 仍由 delayPolicy 统一计算，
 * 不在这里手写超时结果——避免 mock 与真实口径漂移。
 */

const now = Date.now();
const iso = (offsetMin: number) => new Date(now + offsetMin * 60_000).toISOString();
const ptrIso = (offsetMin: number | null) => (offsetMin === null ? null : iso(offsetMin));

export const mockDelayEvents: DelayEvent[] = [
  { id: 1, turnaround_id: 1, flight_no: "CA1234", delay_type: "WEATHER", minutes: 20, root_cause: "本场雷雨", responsibility_team: "OPS", created_at: iso(-70), resolved_at: null, closed: false, open_delay_minutes: 30 },
  { id: 2, turnaround_id: 1, flight_no: "CA1234", delay_type: "CATERING", minutes: 10, root_cause: "餐车晚到", responsibility_team: "CAT", created_at: iso(-50), resolved_at: null, closed: false, open_delay_minutes: 30 },
  { id: 3, turnaround_id: 1, flight_no: "CA1234", delay_type: "BAGGAGE", minutes: 40, root_cause: "传送带故障（已修复）", responsibility_team: "BAG", created_at: iso(-80), resolved_at: ptrIso(-45), closed: true, open_delay_minutes: 30 },
  { id: 4, turnaround_id: 2, flight_no: "CA5678", delay_type: "ATC_FLOW", minutes: 25, root_cause: "流量控制", responsibility_team: "ATC", created_at: iso(-15), resolved_at: null, closed: false, open_delay_minutes: 25 },
  { id: 5, turnaround_id: 3, flight_no: "CA9012", delay_type: "REFUEL", minutes: 15, root_cause: "加油单核对（已闭环）", responsibility_team: "FUL", created_at: iso(-140), resolved_at: ptrIso(-120), closed: true, open_delay_minutes: 0 }
];

const rawTasks: GroundTask[] = [
  { id: 1, turnaround_id: 1, flight_no: "CA1234", task_type: "CATERING", task_type_text: GroundTaskTypeText.CATERING, team_id: "CAT", planned_start: ptrIso(-60), deadline: iso(-15), effective_deadline: iso(-15), open_delay_minutes: 30, open_delay_count: 2, accepted_at: null, actual_finish: null, status: "PENDING", status_text: GroundTaskStatusText.PENDING, blocker_note: "", overdue: false, overdue_minutes: 0, compare_basis: "NOW" },
  { id: 2, turnaround_id: 1, flight_no: "CA1234", task_type: "BAGGAGE", task_type_text: GroundTaskTypeText.BAGGAGE, team_id: "BAG", planned_start: ptrIso(-80), deadline: iso(-45), effective_deadline: iso(-45), open_delay_minutes: 30, open_delay_count: 2, accepted_at: ptrIso(-75), actual_finish: ptrIso(-40), status: "COMPLETED", status_text: GroundTaskStatusText.COMPLETED, blocker_note: "", overdue: false, overdue_minutes: 0, compare_basis: "ACTUAL_FINISH" },
  { id: 3, turnaround_id: 1, flight_no: "CA1234", task_type: "CLEANING", task_type_text: GroundTaskTypeText.CLEANING, team_id: "CLN", planned_start: ptrIso(-85), deadline: iso(-50), effective_deadline: iso(-50), open_delay_minutes: 30, open_delay_count: 2, accepted_at: ptrIso(-80), actual_finish: ptrIso(-10), status: "COMPLETED", status_text: GroundTaskStatusText.COMPLETED, blocker_note: "", overdue: false, overdue_minutes: 0, compare_basis: "ACTUAL_FINISH" },
  { id: 4, turnaround_id: 1, flight_no: "CA1234", task_type: "REFUEL", task_type_text: GroundTaskTypeText.REFUEL, team_id: "FUL", planned_start: ptrIso(-30), deadline: iso(-5), effective_deadline: iso(-5), open_delay_minutes: 30, open_delay_count: 2, accepted_at: null, actual_finish: null, status: "PENDING", status_text: GroundTaskStatusText.PENDING, blocker_note: "", overdue: false, overdue_minutes: 0, compare_basis: "NOW" },
  { id: 5, turnaround_id: 1, flight_no: "CA1234", task_type: "PUSHBACK", task_type_text: GroundTaskTypeText.PUSHBACK, team_id: "PUB", planned_start: ptrIso(-45), deadline: iso(-40), effective_deadline: iso(-40), open_delay_minutes: 30, open_delay_count: 2, accepted_at: ptrIso(-42), actual_finish: null, status: "IN_PROGRESS", status_text: GroundTaskStatusText.IN_PROGRESS, blocker_note: "", overdue: false, overdue_minutes: 0, compare_basis: "NOW" },
  { id: 6, turnaround_id: 2, flight_no: "CA5678", task_type: "CATERING", task_type_text: GroundTaskTypeText.CATERING, team_id: "CAT", planned_start: ptrIso(-10), deadline: iso(10), effective_deadline: iso(10), open_delay_minutes: 25, open_delay_count: 1, accepted_at: null, actual_finish: null, status: "PENDING", status_text: GroundTaskStatusText.PENDING, blocker_note: "", overdue: false, overdue_minutes: 0, compare_basis: "NOW" },
  { id: 7, turnaround_id: 2, flight_no: "CA5678", task_type: "WATER_SERVICE", task_type_text: GroundTaskTypeText.WATER_SERVICE, team_id: "WTR", planned_start: ptrIso(-25), deadline: iso(-30), effective_deadline: iso(-30), open_delay_minutes: 25, open_delay_count: 1, accepted_at: ptrIso(-25), actual_finish: null, status: "IN_PROGRESS", status_text: GroundTaskStatusText.IN_PROGRESS, blocker_note: "", overdue: false, overdue_minutes: 0, compare_basis: "NOW" },
  { id: 8, turnaround_id: 3, flight_no: "CA9012", task_type: "BAGGAGE", task_type_text: GroundTaskTypeText.BAGGAGE, team_id: "BAG", planned_start: ptrIso(-150), deadline: iso(-120), effective_deadline: iso(-120), open_delay_minutes: 0, open_delay_count: 0, accepted_at: ptrIso(-148), actual_finish: ptrIso(-130), status: "COMPLETED", status_text: GroundTaskStatusText.COMPLETED, blocker_note: "", overdue: false, overdue_minutes: 0, compare_basis: "ACTUAL_FINISH" },
  { id: 9, turnaround_id: 3, flight_no: "CA9012", task_type: "REFUEL", task_type_text: GroundTaskTypeText.REFUEL, team_id: "FUL", planned_start: ptrIso(-130), deadline: iso(-100), effective_deadline: iso(-100), open_delay_minutes: 0, open_delay_count: 0, accepted_at: ptrIso(-128), actual_finish: ptrIso(-90), status: "COMPLETED", status_text: GroundTaskStatusText.COMPLETED, blocker_note: "", overdue: false, overdue_minutes: 0, compare_basis: "ACTUAL_FINISH" }
];

// 统一用 delayPolicy 重算生效截止与超时（与看板/列表同一函数）。
export const mockTasks: GroundTask[] = decorateTasks(rawTasks, new Date());

// 航班上的任务总数/完成数/超时数一律从统一口径的任务派生，不手写——避免“条数与分钟对不上”。
function withTaskCounts(base: Omit<FlightTurnaround, "task_total" | "task_completed" | "task_overdue">): FlightTurnaround {
  const list = mockTasks.filter((t) => t.turnaround_id === base.id);
  return {
    ...base,
    task_total: list.length,
    task_completed: list.filter((t) => t.status === "COMPLETED").length,
    task_overdue: list.filter((t) => t.overdue).length
  };
}

export const mockFlights: FlightTurnaround[] = [
  withTaskCounts({
    id: 1, flight_no: "CA1234", aircraft_reg: "B-5861", stand_no: "201",
    arrival_time: iso(-90), departure_time: iso(30), turnaround_status: "IN_SERVICE",
    status_text: TurnaroundStatusText.IN_SERVICE, delay_reason: "天气+餐食延误",
    open_delay_minutes: 30, open_delay_count: 2
  }),
  withTaskCounts({
    id: 2, flight_no: "CA5678", aircraft_reg: "B-6207", stand_no: "118",
    arrival_time: iso(-20), departure_time: iso(40), turnaround_status: "ON_STAND",
    status_text: TurnaroundStatusText.ON_STAND, delay_reason: "",
    open_delay_minutes: 25, open_delay_count: 1
  }),
  withTaskCounts({
    id: 3, flight_no: "CA9012", aircraft_reg: "B-3099", stand_no: "305",
    arrival_time: iso(-160), departure_time: iso(-60), turnaround_status: "DEPARTED",
    status_text: TurnaroundStatusText.DEPARTED, delay_reason: "",
    open_delay_minutes: 0, open_delay_count: 0
  })
];

export const mockResources: GroundResource[] = [
  { id: 1, resource_code: "GPU-01", resource_type: "POWER", location: "201", availability_status: "BOOKED", status_text: "已预约", maintenance_due_at: null, owner_team: "EQP" },
  { id: 2, resource_code: "CAT-07", resource_type: "CATERING", location: "厨房-2", availability_status: "AVAILABLE", status_text: "可用", maintenance_due_at: null, owner_team: "CAT" },
  { id: 3, resource_code: "FUEL-03", resource_type: "REFUEL", location: "油站-A", availability_status: "MAINTENANCE", status_text: "检修中", maintenance_due_at: ptrIso(240), owner_team: "FUL" },
  { id: 4, resource_code: "BLT-12", resource_type: "BAGGAGE", location: "分拣-1", availability_status: "OFFLINE", status_text: "下线", maintenance_due_at: null, owner_team: "BAG" }
];

export const mockBookings: ResourceBooking[] = [
  { id: 1, resource_id: 1, resource_code: "GPU-01", turnaround_id: 1, task_id: 1, start_time: iso(-60), end_time: iso(0), booking_status: "CONFIRMED", status_text: "已确认", conflict_reason: "" },
  { id: 2, resource_id: 1, resource_code: "GPU-01", turnaround_id: 2, task_id: null, start_time: iso(-30), end_time: iso(20), booking_status: "CONFLICT", status_text: "冲突", conflict_reason: "与 GPU-01 既有预约时间重叠" },
  { id: 3, resource_id: 2, resource_code: "CAT-07", turnaround_id: 2, task_id: 6, start_time: iso(-10), end_time: iso(20), booking_status: "HELD", status_text: "占用", conflict_reason: "" },
  { id: 4, resource_id: 3, resource_code: "FUEL-03", turnaround_id: 1, task_id: null, start_time: iso(-40), end_time: iso(-10), booking_status: "RELEASED", status_text: "已释放", conflict_reason: "" }
];

/** 兼容旧引用：mockData 结构。 */
export const mockData = {
  flightTurnaround: mockFlights,
  groundTask: mockTasks,
  groundResource: mockResources,
  resourceBooking: mockBookings,
  delayEvent: mockDelayEvents
};
