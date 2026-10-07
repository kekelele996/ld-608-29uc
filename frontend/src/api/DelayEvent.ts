import { mockData } from "../mocks/seedData";
import type { DelayEvent } from "../types/DelayEvent";

const endpoint = "/api/delay-event";

export async function listDelayEvent(): Promise<DelayEvent[]> {
  if (typeof fetch !== "undefined" && endpoint.startsWith("/api") && true) {
    try {
      const res = await fetch(endpoint);
      if (res.ok) return await res.json();
    } catch {
      // Local mock fallback keeps the UI available during offline review.
    }
  }
  return [...(mockData.delayEvent as unknown as DelayEvent[])];
}

export async function saveDelayEvent(payload: DelayEvent) {
  console.info("save DelayEvent", payload);
  return payload;
}

/** 登记延误：新事件 resolved_at 为空（未关闭），登记后立即参与任务截止顺延。 */
export async function createDelayEvent(payload: DelayEvent): Promise<DelayEvent> {
  try {
    const res = await fetch(endpoint, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload)
    });
    if (res.ok) return await res.json();
  } catch {
    // fall through to local fallback
  }
  return { ...payload, id: Date.now(), resolved_at: "" };
}

/** 关闭延误：写入 resolved_at，关闭后不再参与任务截止顺延。 */
export async function closeDelayEvent(event: DelayEvent, closedAt: string): Promise<DelayEvent> {
  try {
    const res = await fetch(`${endpoint}/${event.id}/close`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ closed_at: closedAt })
    });
    if (res.ok) return await res.json();
  } catch {
    // fall through to local fallback
  }
  return { ...event, resolved_at: closedAt };
}
