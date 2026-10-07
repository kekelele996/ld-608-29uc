import type { GroundTask } from "../types/GroundTask";

export const createDefaultGroundTask = (overrides: Partial<GroundTask> = {}): GroundTask => ({
  id: 0,
  turnaround_id: 0,
  task_type: "CLEANING",
  team_id: 1,
  planned_start: "",
  deadline: "",
  actual_finish: "",
  status: "PENDING",
  blocker_note: "",
  ...overrides
});

export const createGroundTaskForm = createDefaultGroundTask;
export const createGroundTaskResponse = createDefaultGroundTask;
