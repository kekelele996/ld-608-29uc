import { create } from "zustand";
import { listGroundTask, signOffGroundTask } from "../api/GroundTask";
import { TASK_LOG_SIGN_OFF } from "../constants/logTemplates";
import type { GroundTask } from "../types/GroundTask";

type State = {
  rows: GroundTask[];
  loading: boolean;
  load: () => Promise<void>;
  signOff: (id: number) => Promise<void>;
};

export const useGroundTaskStore = create<State>((set, get) => ({
  rows: [],
  loading: false,
  async load() {
    set({ loading: true });
    set({ rows: await listGroundTask(), loading: false });
  },
  async signOff(id) {
    const task = get().rows.find((row) => row.id === id);
    if (!task || task.actual_finish) return;
    const signedAt = new Date().toISOString();
    const updated = await signOffGroundTask(task, signedAt);
    console.info(TASK_LOG_SIGN_OFF, { taskId: id, signedAt });
    set({ rows: get().rows.map((row) => (row.id === id ? updated : row)) });
  }
}));
