import { create } from "zustand";
import { closeDelayEvent, createDelayEvent, listDelayEvent } from "../api/DelayEvent";
import { DELAY_LOG_CLOSE, DELAY_LOG_REGISTER } from "../constants/logTemplates";
import type { DelayEvent } from "../types/DelayEvent";

type State = {
  rows: DelayEvent[];
  loading: boolean;
  load: () => Promise<void>;
  register: (payload: DelayEvent) => Promise<void>;
  close: (id: number) => Promise<void>;
};

export const useDelayEventStore = create<State>((set, get) => ({
  rows: [],
  loading: false,
  async load() {
    set({ loading: true });
    set({ rows: await listDelayEvent(), loading: false });
  },
  async register(payload) {
    const created = await createDelayEvent({ ...payload, resolved_at: "" });
    console.info(DELAY_LOG_REGISTER, created);
    set({ rows: [...get().rows, created] });
  },
  async close(id) {
    const event = get().rows.find((row) => row.id === id);
    if (!event || event.resolved_at) return;
    const closedAt = new Date().toISOString();
    const updated = await closeDelayEvent(event, closedAt);
    console.info(DELAY_LOG_CLOSE, { delayId: id, closedAt });
    set({ rows: get().rows.map((row) => (row.id === id ? updated : row)) });
  }
}));
