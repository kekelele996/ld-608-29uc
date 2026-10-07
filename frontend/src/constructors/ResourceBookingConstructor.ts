import type { ResourceBooking } from "../types/ResourceBooking";

export const createDefaultResourceBooking = (overrides: Partial<ResourceBooking> = {}): ResourceBooking => ({
  id: 0,
  resource_id: 0,
  resource_code: "",
  turnaround_id: 0,
  task_id: null,
  start_time: new Date().toISOString(),
  end_time: new Date().toISOString(),
  booking_status: "HELD",
  status_text: "占用",
  conflict_reason: "",
  ...overrides
});

export const createResourceBookingForm = createDefaultResourceBooking;
export const createResourceBookingResponse = createDefaultResourceBooking;
