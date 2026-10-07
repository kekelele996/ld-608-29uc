import { useEffect, useMemo, useState } from "react";
import { useDelayEventStore } from "../stores/DelayEventStore";
import { useFlightTurnaroundStore } from "../stores/FlightTurnaroundStore";
import { useSessionStore } from "../stores/SessionStore";
import { DelayTag } from "../components/common/DelayTag";
import { StatusBadge } from "../components/common/StatusBadge";
import { EmptyState } from "../components/common/EmptyState";
import { DeadlinePolicyBanner } from "../components/common/DeadlinePolicyBanner";
import { DelayType, DelayTypeText } from "../constants/GroundTaskStatus";
import { ERROR_MESSAGES } from "../constants/errorMessages";
import type { ApiError } from "../api/http";
import { formatDate } from "../utils/formatters";

const WRITE_ROLES = ["DISPATCHER", "SUPERVISOR"];

export function DelaysPage() {
  const { rows, loading, error, load, register, resolve } = useDelayEventStore();
  const flights = useFlightTurnaroundStore((s) => s.rows);
  const loadFlights = useFlightTurnaroundStore((s) => s.load);
  const role = useSessionStore((s) => s.user?.role ?? "");
  const canWrite = WRITE_ROLES.includes(role);

  const [turnaroundId, setTurnaroundId] = useState<number>(0);
  const [delayType, setDelayType] = useState<string>("WEATHER");
  const [minutes, setMinutes] = useState<number>(15);
  const [team, setTeam] = useState("");
  const [cause, setCause] = useState("");
  const [formError, setFormError] = useState<string | null>(null);
  const [notice, setNotice] = useState<string | null>(null);

  useEffect(() => {
    void load();
    void loadFlights();
  }, [load, loadFlights]);
  useEffect(() => {
    if (!turnaroundId && flights[0]) setTurnaroundId(flights[0].id);
  }, [flights, turnaroundId]);

  const openTotal = useMemo(() => rows.filter((r) => !r.resolved_at).reduce((s, r) => s + r.minutes, 0), [rows]);

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    setFormError(null);
    if (!turnaroundId) {
      setFormError("请选择航班");
      return;
    }
    if (minutes <= 0) {
      setFormError(ERROR_MESSAGES.DEADLINE_NOT_EXTENDED);
      return;
    }
    try {
      await register({ turnaround_id: turnaroundId, delay_type: delayType, minutes, root_cause: cause, responsibility_team: team });
      setNotice(`已登记 ${minutes} 分钟延误，该航班未完成任务的生效截止已顺延 ${minutes} 分钟（原计划截止不变）`);
      setCause("");
      setTimeout(() => setNotice(null), 4000);
    } catch (err) {
      const code = (err as ApiError).code;
      setFormError(ERROR_MESSAGES[code as keyof typeof ERROR_MESSAGES] ?? (err as Error).message);
    }
  }

  async function close(id: number) {
    setFormError(null);
    try {
      await resolve(id);
      setNotice("延误已关闭，该笔分钟不再顺延任务截止（超时按剩余未关闭延误重算）");
      setTimeout(() => setNotice(null), 4000);
    } catch (err) {
      const code = (err as ApiError).code;
      setFormError(ERROR_MESSAGES[code as keyof typeof ERROR_MESSAGES] ?? (err as Error).message);
    }
  }

  return (
    <section className="page">
      <header className="page-head">
        <div>
          <p className="eyebrow">ground-turn / delays</p>
          <h1>延误归因</h1>
          <p className="sub">未关闭延误累计 <b className="warn-text">{openTotal} 分钟</b>，正顺延相关任务的生效截止。</p>
        </div>
      </header>

      <DeadlinePolicyBanner />
      {notice ? <div className="alert ok">{notice}</div> : null}
      {error ? <div className="alert warn">{error}</div> : null}
      {formError ? <div className="alert danger">{formError}</div> : null}

      {canWrite ? (
        <form className="panel delay-form" onSubmit={submit}>
          <h2>登记延误（登记后立即顺延，不改任务原计划截止）</h2>
          <div className="form-grid">
            <label>航班
              <select value={turnaroundId} onChange={(e) => setTurnaroundId(Number(e.target.value))}>
                {flights.map((f) => <option key={f.id} value={f.id}>{f.flight_no}（机位 {f.stand_no}）</option>)}
              </select>
            </label>
            <label>类型
              <select value={delayType} onChange={(e) => setDelayType(e.target.value)}>
                {DelayType.map((d) => <option key={d} value={d}>{DelayTypeText[d]}</option>)}
              </select>
            </label>
            <label>延误分钟（正数才顺延）
              <input type="number" min={1} value={minutes} onChange={(e) => setMinutes(Number(e.target.value))} />
            </label>
            <label>责任班组
              <input value={team} onChange={(e) => setTeam(e.target.value)} placeholder="如 OPS / CAT / BAG" />
            </label>
            <label className="wide">根因
              <input value={cause} onChange={(e) => setCause(e.target.value)} placeholder="延误根因说明" />
            </label>
          </div>
          <button className="btn primary" type="submit">登记延误</button>
        </form>
      ) : <p className="dim">当前角色仅可查看；登记/关闭延误需地勤调度或运行督导。</p>}

      <div className="panel">
        <h2>延误事件清单</h2>
        {loading && rows.length === 0 ? <EmptyState title="加载中…" /> : rows.length === 0 ? <EmptyState title="暂无延误" /> : (
          <div className="table-wrap">
            <table className="data-table">
              <thead><tr><th>#</th><th>航班</th><th>类型</th><th>分钟</th><th>责任</th><th>根因</th><th>登记时间</th><th>关闭时间</th><th>状态/顺延</th><th>操作</th></tr></thead>
              <tbody>
                {rows.map((d) => (
                  <tr key={d.id} className={d.resolved_at ? "row-closed" : ""}>
                    <td>{d.id}</td>
                    <td>{d.flight_no}</td>
                    <td>{DelayTypeText[d.delay_type as keyof typeof DelayTypeText] ?? d.delay_type}</td>
                    <td>{d.minutes}</td>
                    <td>{d.responsibility_team || "—"}</td>
                    <td>{d.root_cause || "—"}</td>
                    <td>{formatDate(d.created_at)}</td>
                    <td>{d.resolved_at ? formatDate(d.resolved_at) : <em className="dim">未关闭</em>}</td>
                    <td>{d.resolved_at
                      ? <StatusBadge value="CLOSED" tone="muted" text="已关闭·不顺延" />
                      : <DelayTag minutes={d.minutes} label="顺延" />}</td>
                    <td>{canWrite && !d.resolved_at
                      ? <button className="btn small" onClick={() => void close(d.id)}>关闭（停止顺延）</button>
                      : <span className="dim">—</span>}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </section>
  );
}
