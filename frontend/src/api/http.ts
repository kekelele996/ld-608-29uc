import { useSessionStore } from "../stores/SessionStore";

/**
 * 统一请求封装：前缀 /api（nginx 反代到后端，禁止硬编码 localhost），
 * 自动带 JWT、解包 {ok,data}、401 时清空会话。
 */
const BASE = "/api";

export class ApiError extends Error {
  code: string;
  status: number;
  constructor(status: number, code: string, message: string) {
    super(message);
    this.status = status;
    this.code = code;
  }
}

type Options = {
  method?: string;
  body?: unknown;
  auth?: boolean;
};

export async function request<T>(path: string, options: Options = {}): Promise<T> {
  const { method = "GET", body, auth = true } = options;
  const headers: Record<string, string> = { "Content-Type": "application/json" };
  if (auth) {
    const token = useSessionStore.getState().token;
    if (token) headers.Authorization = `Bearer ${token}`;
  }
  let res: Response;
  try {
    res = await fetch(BASE + path, {
      method,
      headers,
      body: body !== undefined ? JSON.stringify(body) : undefined
    });
  } catch {
    throw new ApiError(0, "NETWORK_OFFLINE", "后端不可达，已切换本地演示数据");
  }
  const payload = await res.json().catch(() => null);
  if (!res.ok) {
    const err = payload?.error ?? {};
    if (res.status === 401) {
      useSessionStore.getState().clear();
    }
    throw new ApiError(res.status, err.code ?? "INTERNAL_ERROR", err.message ?? "请求失败");
  }
  return (payload?.data ?? payload) as T;
}
