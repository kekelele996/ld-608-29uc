import { create } from "zustand";
import { persist } from "zustand/middleware";
import type { SessionUser } from "../types/System";

interface SessionState {
  user: SessionUser | null;
  token: string | null;
  setSession: (user: SessionUser) => void;
  clear: () => void;
}

export const useSessionStore = create<SessionState>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      setSession: (user) => set({ user, token: user.token }),
      clear: () => set({ user: null, token: null })
    }),
    { name: "ground-turn-session" }
  )
);
