import { useState } from "react";
import { login } from "../../api/System";
import { useSessionStore } from "../../stores/SessionStore";
import { ERROR_MESSAGES } from "../../constants/errorMessages";
import type { ApiError } from "../../api/http";

const DEMO = [
  { username: "dispatcher", role: "地勤调度" },
  { username: "team", role: "班组" },
  { username: "resource", role: "资源管理员" },
  { username: "supervisor", role: "运行督导" }
];

export function LoginModal({ onClose }: { onClose: () => void }) {
  const setSession = useSessionStore((s) => s.setSession);
  const [username, setUsername] = useState("dispatcher");
  const [password, setPassword] = useState("demo123");
  const [error, setError] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    setBusy(true);
    setError(null);
    try {
      const user = await login({ username, password });
      setSession(user);
      onClose();
    } catch (err) {
      const code = (err as ApiError).code;
      setError(ERROR_MESSAGES[code as keyof typeof ERROR_MESSAGES] ?? (err as Error).message);
    } finally {
      setBusy(false);
    }
  }

  return (
    <div className="modal-mask">
      <form className="modal" onSubmit={submit}>
        <h2>登录 · 航空地勤周转保障平台</h2>
        <label>账号<input value={username} onChange={(e) => setUsername(e.target.value)} /></label>
        <label>口令<input type="password" value={password} onChange={(e) => setPassword(e.target.value)} /></label>
        {error ? <p className="form-error">{error}</p> : null}
        <button className="btn primary" disabled={busy} type="submit">{busy ? "登录中…" : "登录"}</button>
        <div className="demo-users">
          <span>演示账号（口令均为 demo123，可切换角色验证按钮显隐）：</span>
          <div className="demo-chips">
            {DEMO.map((d) => (
              <button type="button" key={d.username} className="chip" onClick={() => setUsername(d.username)}>
                {d.role} · {d.username}
              </button>
            ))}
          </div>
        </div>
      </form>
    </div>
  );
}
