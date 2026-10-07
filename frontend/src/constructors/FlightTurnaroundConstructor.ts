import type { FlightTurnaround } from "../types/FlightTurnaround";

export const createDefaultFlightTurnaround = (overrides: Partial<FlightTurnaround> = {}): FlightTurnaround => ({
  id: 0,
  flight_no: "",
  aircraft_reg: "",
  stand_no: "",
  arrival_time: new Date().toISOString(),
  departure_time: new Date().toISOString(),
  turnaround_status: "ARRIVING",
  status_text: "进近中",
  delay_reason: "",
  open_delay_minutes: 0,
  open_delay_count: 0,
  task_total: 0,
  task_completed: 0,
  task_overdue: 0,
  ...overrides
});

export const createFlightTurnaroundForm = createDefaultFlightTurnaround;
export const createFlightTurnaroundResponse = createDefaultFlightTurnaround;
