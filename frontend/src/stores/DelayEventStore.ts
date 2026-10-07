import { create } from "zustand";
import { listDelayEvents, registerDelay, resolveDelay } from "../api/DelayEvent";
import type { DelayEvent, DelayRegisterPayload } from "../types/DelayEvent";

interface DelayEventState {
  rows: DelayEvent[];
  loading: boolean;
  error: string | null;
  load: () => Promise<void>;
  register: (payload: DelayRegisterPayload) => Promise<void>;
  resolve: (id: number) => Promise<void>;
  clearError: () => void;
}

export const useDelayEventStore = create<DelayEventState>((set) => ({
  rows: [],
  loading: false,
  error: null,
  async load() {
    set({ loading: true, error: null });
    try {
      set({ rows: await listDelayEvents() });
    } catch (e) {
      set({ error: (e as Error).message });
    } finally {
      set({ loading: false });
    }
  },
  async register(payload) {
    set({ error: null });
    await registerDelay(payload);
    await useDelayEventStore.getState().load();
  },
  async resolve(id) {
    set({ error: null });
    await resolveDelay(id);
    await useDelayEventStore.getState().load();
  },
  clearError: () => set({ error: null })
}));
