import { useEffect, useMemo, useState } from "react";
import { useGroundTaskStore } from "../stores/GroundTaskStore";
import { useFlightTurnaroundStore } from "../stores/FlightTurnaroundStore";
import { useSessionStore } from "../stores/SessionStore";
import { GroundTaskTable } from "../components/common/GroundTaskTable";
import { DeadlinePolicyBanner } from "../components/common/DeadlinePolicyBanner";
import { EmptyState } from "../components/common/EmptyState";
import { GroundTaskStatus, GroundTaskStatusText } from "../constants/GroundTaskStatus";
import { ERROR_MESSAGES } from "../constants/errorMessages";
import type { ApiError } from "../api/http";
import type { GroundTask } from "../types/GroundTask";

const SIGN_ROLES = ["TEAM", "DISPATCHER", "SUPERVISOR"];

export function TasksPage() {
  const { rows, loading, error, load, accept, complete } = useGroundTaskStore();
  const loadFlights = useFlightTurnaroundStore((s) => s.load);
  const role = useSessionStore((s) => s.user?.role ?? "");
  const canSign = SIGN_ROLES.includes(role);

  const [statusFilter, setStatusFilter] = useState<string>("ALL");
  const [flightFilter, setFlightFilter] = useState<string>("ALL");
  const [overdueOnly, setOverdueOnly] = useState(false);
  const [actionError, setActionError] = useState<string | null>(null);
  const [notice, setNotice] = useState<string | null>(null);

  useEffect(() => {
    void load();
    void loadFlights();
  }, [load, loadFlights]);

  const flights = useFlightTurnaroundStore((s) => s.rows);

  const filtered = useMemo(() => {
    return rows.filter((t) => {
      if (statusFilter !== "ALL" && t.status !== statusFilter) return false;
      if (flightFilter !== "ALL" && String(t.turnaround_id) !== flightFilter) return false;
      if (overdueOnly && !t.overdue) return false;
      return true;
    });
  }, [rows, statusFilter, flightFilter, overdueOnly]);

  async function guard(fn: () => Promise<void>, ok: string) {
    setActionError(null);
    try {
      await fn();
      setNotice(ok);
      setTimeout(() => setNotice(null), 2500);
    } catch (e) {
      const code = (e as ApiError).code;
      setActionError(ERROR_MESSAGES[code as keyof typeof ERROR_MESSAGES] ?? (e as Error).message);
    }
  }

  const pendingCount = rows.filter((t: GroundTask) => t.accepted_at === null).length;

  return (
    <section className="page">
      <header className="page-head">
        <div>
          <p className="eyebrow">ground-turn / tasks</p>
          <h1>地勤任务</h1>
          <p className="sub">共 {rows.length} 项，未签收 <b>{pendingCount}</b> 项，超时 <b className="danger-text">{rows.filter((t) => t.overdue).length}</b> 项（按生效截止）。</p>
        </div>
      </header>

      <DeadlinePolicyBanner />

      {notice ? <div className="alert ok">{notice}</div> : null}
      {error ? <div className="alert warn">列表加载失败：{error}（当前可能为离线演示数据）</div> : null}
      {actionError ? <div className="alert danger">{actionError}</div> : null}

      <div className="filters">
        <label>状态
          <select value={statusFilter} onChange={(e) => setStatusFilter(e.target.value)}>
            <option value="ALL">全部</option>
            {GroundTaskStatus.map((s) => <option key={s} value={s}>{GroundTaskStatusText[s]}</option>)}
          </select>
        </label>
        <label>航班
          <select value={flightFilter} onChange={(e) => setFlightFilter(e.target.value)}>
            <option value="ALL">全部</option>
            {flights.map((f) => <option key={f.id} value={String(f.id)}>{f.flight_no}</option>)}
          </select>
        </label>
        <label className="check">
          <input type="checkbox" checked={overdueOnly} onChange={(e) => setOverdueOnly(e.target.checked)} /> 只看超时
        </label>
        {!canSign && role ? <span className="dim">当前角色无签收/完成权限（需班组/调度/督导）</span> : null}
      </div>

      {loading && rows.length === 0 ? <EmptyState title="加载中…" /> : filtered.length === 0
        ? <EmptyState title="没有符合条件的任务" hint="调整状态/航班/超时筛选再试" />
        : <GroundTaskTable
            tasks={filtered}
            canSign={canSign}
            onAccept={(id) => void guard(() => accept(id), "任务已签收，签收时间已记录")}
            onComplete={(id) => void guard(() => complete(id), "任务已完成，按实际完成时间完成超时判定")}
          />}
    </section>
  );
}
