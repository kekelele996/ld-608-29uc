import { request } from "./http";
import type { Dashboard, LoginPayload, SessionUser, AuditLog } from "../types/System";
import { mockFlights, mockTasks, mockDelayEvents } from "../mocks/seedData";
import { decorateTasks } from "../utils/delayPolicy";

export async function login(payload: LoginPayload): Promise<SessionUser> {
  return request<SessionUser>("/auth/login", { method: "POST", body: payload, auth: false });
}

export async function getDashboard(): Promise<Dashboard> {
  try {
    return await request<Dashboard>("/dashboard");
  } catch {
    // 离线兜底：用与线上相同的装配方式现算，保证看板口径不漂移。
    const tasks = decorateTasks(structuredClone(mockTasks));
    const overdueTasks = tasks.filter((t) => t.overdue);
    const openDelay = mockDelayEvents.filter((d) => !d.resolved_at);
    return {
      turnaround_total: mockFlights.length,
      active_turnaround: mockFlights.filter((f) => f.turnaround_status !== "DEPARTED").length,
      open_delay_minutes: openDelay.reduce((s, d) => s + d.minutes, 0),
      open_delay_events: openDelay.length,
      overdue_task_count: overdueTasks.length,
      task_total: tasks.length,
      task_accepted: tasks.filter((t) => t.accepted_at).length,
      task_completed: tasks.filter((t) => t.status === "COMPLETED").length,
      overdue_tasks: overdueTasks,
      turnarounds: mockFlights
    };
  }
}

export async function listAuditLogs(): Promise<AuditLog[]> {
  try {
    return await request<AuditLog[]>("/audit-logs");
  } catch {
    return [];
  }
}
