import type { GroundTask } from "./GroundTask";

export interface FlightTurnaround {
  id: number;
  flight_no: string;
  aircraft_reg: string;
  stand_no: string;
  arrival_time: string;
  departure_time: string;
  turnaround_status: string;
  status_text: string;
  delay_reason: string;
  open_delay_minutes: number;
  open_delay_count: number;
  task_total: number;
  task_completed: number;
  task_overdue: number;
  tasks?: GroundTask[];
}
