import type { GroundTask } from "./GroundTask";
import type { FlightTurnaround } from "./FlightTurnaround";

export interface Dashboard {
  turnaround_total: number;
  active_turnaround: number;
  open_delay_minutes: number;
  open_delay_events: number;
  overdue_task_count: number;
  task_total: number;
  task_accepted: number;
  task_completed: number;
  overdue_tasks: GroundTask[];
  turnarounds: FlightTurnaround[];
}

export interface LoginPayload {
  username: string;
  password: string;
}

export interface SessionUser {
  token: string;
  username: string;
  role: string;
  role_text: string;
  team_code: string;
  display_name: string;
}

export interface AuditLog {
  id: number;
  actor: string;
  action: string;
  target_type: string;
  target_id: string;
  detail: string;
  created_at: string;
}
