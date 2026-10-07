import { request } from "./http";
import { mockTasks } from "../mocks/seedData";
import type { GroundTask } from "../types/GroundTask";

export async function listGroundTasks(): Promise<GroundTask[]> {
  try {
    return await request<GroundTask[]>("/ground-tasks");
  } catch {
    return structuredClone(mockTasks);
  }
}

export async function acceptGroundTask(id: number): Promise<GroundTask> {
  return request<GroundTask>(`/ground-tasks/${id}/accept`, { method: "POST" });
}

export async function completeGroundTask(id: number): Promise<GroundTask> {
  return request<GroundTask>(`/ground-tasks/${id}/complete`, { method: "POST" });
}

export async function blockGroundTask(id: number, blockerNote: string): Promise<GroundTask> {
  return request<GroundTask>(`/ground-tasks/${id}/block`, { method: "POST", body: { blocker_note: blockerNote } });
}

export interface DispatchPayload {
  turnaround_id: number;
  task_type: string;
  team_id?: string;
  planned_start?: string | null;
  deadline: string;
}

export async function dispatchGroundTask(payload: DispatchPayload): Promise<GroundTask> {
  return request<GroundTask>("/ground-tasks", { method: "POST", body: payload });
}
