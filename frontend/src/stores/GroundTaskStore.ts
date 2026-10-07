import { create } from "zustand";
import {
  listGroundTasks,
  acceptGroundTask,
  completeGroundTask,
  blockGroundTask,
  dispatchGroundTask,
  type DispatchPayload
} from "../api/GroundTask";
import type { GroundTask } from "../types/GroundTask";

interface GroundTaskState {
  rows: GroundTask[];
  loading: boolean;
  error: string | null;
  load: () => Promise<void>;
  accept: (id: number) => Promise<void>;
  complete: (id: number) => Promise<void>;
  block: (id: number, note: string) => Promise<void>;
  dispatch: (payload: DispatchPayload) => Promise<void>;
  clearError: () => void;
}

export const useGroundTaskStore = create<GroundTaskState>((set) => ({
  rows: [],
  loading: false,
  error: null,
  async load() {
    set({ loading: true, error: null });
    try {
      set({ rows: await listGroundTasks() });
    } catch (e) {
      set({ error: (e as Error).message });
    } finally {
      set({ loading: false });
    }
  },
  async accept(id) {
    set({ error: null });
    await acceptGroundTask(id);
    await useGroundTaskStore.getState().load();
  },
  async complete(id) {
    set({ error: null });
    await completeGroundTask(id);
    await useGroundTaskStore.getState().load();
  },
  async block(id, note) {
    set({ error: null });
    await blockGroundTask(id, note);
    await useGroundTaskStore.getState().load();
  },
  async dispatch(payload) {
    set({ error: null });
    await dispatchGroundTask(payload);
    await useGroundTaskStore.getState().load();
  },
  clearError: () => set({ error: null })
}));
