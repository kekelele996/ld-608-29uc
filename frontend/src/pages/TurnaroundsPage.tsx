import { useEffect } from "react";
import { useFlightTurnaroundStore } from "../stores/FlightTurnaroundStore";
import { TurnaroundTimeline } from "../components/common/TurnaroundTimeline";
import { EmptyState } from "../components/common/EmptyState";
import { useTurnaroundProgress } from "../hooks/useTurnaroundProgress";
import { StatCard } from "../components/common/StatCard";
import { formatMinutes } from "../utils/formatters";

export function TurnaroundsPage() {
  const { rows, loading, load } = useFlightTurnaroundStore();
  useEffect(() => {
    void load();
  }, [load]);
  const progress = useTurnaroundProgress(rows);

  return (
    <section className="page">
      <header className="page-head">
        <div>
          <p className="eyebrow">ground-turn / turnarounds</p>
          <h1>航班过站</h1>
          <p className="sub">{rows.length} 个航班，任务完成率 {progress.completionRate}%。</p>
        </div>
      </header>

      <section className="metrics">
        <StatCard label="航班数" value={rows.length} />
        <StatCard label="未关闭顺延" value={formatMinutes(progress.openDelayMinutes)} tone="warn" />
        <StatCard label="超时任务" value={progress.overdue} tone="danger" />
        <StatCard label="完成率" value={`${progress.completionRate}%`} tone="ok" />
      </section>

      {loading && rows.length === 0 ? <EmptyState title="加载中…" /> : (
        <div className="timeline-grid">
          {rows.map((f) => <TurnaroundTimeline key={f.id} flight={f} />)}
        </div>
      )}
    </section>
  );
}
