import { create } from "zustand";
import { getDashboard } from "../api/System";
import type { Dashboard } from "../types/System";

type State = {
  data: Dashboard | null;
  loading: boolean;
  load: () => Promise<void>;
};

export const useDashboardStore = create<State>((set) => ({
  data: null,
  loading: false,
  async load() {
    set({ loading: true });
    try {
      set({ data: await getDashboard(), loading: false });
    } catch {
      set({ loading: false });
    }
  }
}));
