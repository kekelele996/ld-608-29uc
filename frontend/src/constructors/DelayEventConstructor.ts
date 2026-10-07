import type { DelayEvent } from "../types/DelayEvent";

export const createDefaultDelayEvent = (overrides: Partial<DelayEvent> = {}): DelayEvent => ({
  id: 0,
  turnaround_id: 0,
  delay_type: "CATERING",
  minutes: 0,
  root_cause: "",
  responsibility_team: "",
  resolved_at: "",
  ...overrides
});

export const createDelayEventForm = createDefaultDelayEvent;
export const createDelayEventResponse = createDefaultDelayEvent;
