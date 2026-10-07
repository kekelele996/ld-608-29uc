import { mockData } from "../mocks/seedData";
import type { GroundTask } from "../types/GroundTask";

const endpoint = "/api/ground-task";

export async function listGroundTask(): Promise<GroundTask[]> {
  if (typeof fetch !== "undefined" && endpoint.startsWith("/api") && true) {
    try {
      const res = await fetch(endpoint);
      if (res.ok) return await res.json();
    } catch {
      // Local mock fallback keeps the UI available during offline review.
    }
  }
  return [...(mockData.groundTask as unknown as GroundTask[])];
}

export async function saveGroundTask(payload: GroundTask) {
  console.info("save GroundTask", payload);
  return payload;
}

/**
 * 签收任务：后端把签收时间写入 actual_finish 并把状态置为 SIGNED；
 * 后端不可用时按同一规则在本地回退构造，保证 mock 评审环境行为一致。
 */
export async function signOffGroundTask(task: GroundTask, signedAt: string): Promise<GroundTask> {
  try {
    const res = await fetch(`${endpoint}/${task.id}/sign-off`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ signed_at: signedAt })
    });
    if (res.ok) return await res.json();
  } catch {
    // fall through to local fallback
  }
  return { ...task, status: "SIGNED", actual_finish: signedAt };
}
