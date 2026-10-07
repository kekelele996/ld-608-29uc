export interface ResourceBooking {
  id: number;
  resource_id: number;
  resource_code: string;
  turnaround_id: number;
  task_id: number | null;
  start_time: string;
  end_time: string;
  booking_status: string;
  status_text: string;
  conflict_reason: string;
}
