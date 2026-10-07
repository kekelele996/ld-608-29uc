import type { DelayEvent } from "../types/DelayEvent";

/** 延误登记表单默认对象：新登记默认未关闭（resolved_at=null），立即顺延。 */
export const createDefaultDelayEvent = (overrides: Partial<DelayEvent> = {}): DelayEvent => ({
  id: 0,
  turnaround_id: 0,
  flight_no: "",
  delay_type: "WEATHER",
  minutes: 15,
  root_cause: "",
  responsibility_team: "",
  created_at: new Date().toISOString(),
  resolved_at: null,
  closed: false,
  open_delay_minutes: 0,
  ...overrides
});

export const createDelayEventForm = createDefaultDelayEvent;
export const createDelayEventResponse = createDefaultDelayEvent;
