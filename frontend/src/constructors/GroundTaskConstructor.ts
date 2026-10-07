import type { GroundTask, GroundTaskStatus } from "../types/GroundTask";

/** 派工表单默认对象，页面/store 不得散写默认结构。 */
export const createDefaultGroundTask = (overrides: Partial<GroundTask> = {}): GroundTask => ({
  id: 0,
  turnaround_id: 0,
  flight_no: "",
  task_type: "CATERING",
  task_type_text: "",
  team_id: "",
  planned_start: null,
  deadline: new Date().toISOString(),
  effective_deadline: new Date().toISOString(),
  open_delay_minutes: 0,
  open_delay_count: 0,
  accepted_at: null,
  actual_finish: null,
  status: "PENDING" as GroundTaskStatus,
  status_text: "待签收",
  blocker_note: "",
  overdue: false,
  overdue_minutes: 0,
  compare_basis: "NOW",
  ...overrides
});

export const createGroundTaskForm = createDefaultGroundTask;
export const createGroundTaskResponse = createDefaultGroundTask;
