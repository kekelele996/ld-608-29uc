import { useEffect } from "react";
import { useDashboardStore } from "../stores/DashboardStore";
import { StatCard } from "../components/common/StatCard";
import { OverdueTaskCard } from "../components/common/OverdueTaskCard";
import { TurnaroundTimeline } from "../components/common/TurnaroundTimeline";
import { DeadlinePolicyBanner } from "../components/common/DeadlinePolicyBanner";
import { EmptyState } from "../components/common/EmptyState";
import { useTurnaroundProgress } from "../hooks/useTurnaroundProgress";
import { formatMinutes } from "../utils/formatters";

export function DashboardPage() {
  const { data, loading, load } = useDashboardStore();
  useEffect(() => {
    void load();
  }, [load]);
  const progress = useTurnaroundProgress(data?.turnarounds ?? []);

  if (loading && !data) return <section className="page"><EmptyState title="加载看板中…" /></section>;
  if (!data) return <section className="page"><EmptyState title="暂无看板数据" /></section>;

  return (
    <section className="page">
      <header className="page-head">
        <div>
          <p className="eyebrow">ground-turn / dashboard</p>
          <h1>过站运行看板</h1>
          <p className="sub">在保障航班 {data.active_turnaround} / {data.turnaround_total}；任务完成率 {progress.completionRate}%。</p>
        </div>
      </header>

      <section className="metrics">
        <StatCard label="超时任务（按生效截止）" value={data.overdue_task_count} tone="danger" hint="与任务列表同一口径" />
        <StatCard label="未关闭延误累计" value={formatMinutes(data.open_delay_minutes)} tone="warn" hint={`${data.open_delay_events} 起未关闭，关闭过的不计`} />
        <StatCard label="任务总数" value={data.task_total} hint={`已签收 ${data.task_accepted} · 已完成 ${data.task_completed}`} />
        <StatCard label="在保障航班" value={data.active_turnaround} tone="ok" />
      </section>

      <DeadlinePolicyBanner />

      <section className="workbench">
        <div className="panel wide">
          <h2>超时任务 <span className="dim">（签收/完成任务按实际完成时间，其余按当前时间）</span></h2>
          {data.overdue_tasks.length === 0
            ? <EmptyState title="当前没有超时任务" hint="生效截止之内，或已关闭延误不再顺延" />
            : <div className="overdue-list">
                {data.overdue_tasks.map((t) => <OverdueTaskCard key={t.id} task={t} />)}
              </div>}
        </div>
        <div className="panel">
          <h2>未关闭延误 vs 超时对账</h2>
          <ul className="reconcile">
            {data.turnarounds.map((f) => (
              <li key={f.id}>
                <strong>{f.flight_no}</strong>
                <span>未关闭顺延 {formatMinutes(f.open_delay_minutes)}（{f.open_delay_count} 起）</span>
                <span className={f.task_overdue > 0 ? "danger-text" : ""}>超时任务 {f.task_overdue} / {f.task_total}</span>
              </li>
            ))}
          </ul>
          <p className="policy-note">
            对账方式：每条任务的生效截止 = 原计划 + 该航班未关闭延误分钟；超时条数与顺延分钟因此一一对应，不会因关闭过的延误或原计划时间错位。
          </p>
        </div>
      </section>

      <section className="panel">
        <h2>航班进度</h2>
        <div className="timeline-grid">
          {data.turnarounds.map((f) => <TurnaroundTimeline key={f.id} flight={f} />)}
        </div>
      </section>
    </section>
  );
}
