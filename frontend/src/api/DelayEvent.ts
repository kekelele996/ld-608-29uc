import { request } from "./http";
import { mockDelayEvents } from "../mocks/seedData";
import type { DelayEvent, DelayRegisterPayload } from "../types/DelayEvent";

export async function listDelayEvents(): Promise<DelayEvent[]> {
  try {
    return await request<DelayEvent[]>("/delay-events");
  } catch {
    return structuredClone(mockDelayEvents);
  }
}

export async function registerDelay(payload: DelayRegisterPayload): Promise<DelayEvent> {
  return request<DelayEvent>("/delay-events", { method: "POST", body: payload });
}

export async function resolveDelay(id: number): Promise<DelayEvent> {
  return request<DelayEvent>(`/delay-events/${id}/resolve`, { method: "POST" });
}
