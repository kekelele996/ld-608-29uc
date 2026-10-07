import { useEffect } from "react";
import { useGroundResourceStore } from "../stores/GroundResourceStore";
import { useResourceBookingStore } from "../stores/ResourceBookingStore";
import { ResourceCalendar } from "../components/common/ResourceCalendar";
import { StatusBadge } from "../components/common/StatusBadge";
import { EmptyState } from "../components/common/EmptyState";
import { useResourceConflict } from "../hooks/useResourceConflict";
import { formatDate } from "../utils/formatters";

export function ResourcesPage() {
  const resources = useGroundResourceStore((s) => s.rows);
  const loadResources = useGroundResourceStore((s) => s.load);
  const bookings = useResourceBookingStore((s) => s.rows);
  const loadBookings = useResourceBookingStore((s) => s.load);

  useEffect(() => {
    void loadResources();
    void loadBookings();
  }, [loadResources, loadBookings]);

  const { conflictCount } = useResourceConflict(bookings);

  return (
    <section className="page">
      <header className="page-head">
        <div>
          <p className="eyebrow">ground-turn / resources</p>
          <h1>资源调度</h1>
          <p className="sub">资源 {resources.length} 项，预约冲突 <b className={conflictCount ? "danger-text" : ""}>{conflictCount}</b> 起。</p>
        </div>
      </header>

      <section className="panel">
        <h2>资源台账</h2>
        {resources.length === 0 ? <EmptyState title="暂无资源" /> : (
          <div className="table-wrap">
            <table className="data-table">
              <thead><tr><th>编码</th><th>类型</th><th>位置</th><th>状态</th><th>归属班组</th><th>检修到期</th></tr></thead>
              <tbody>
                {resources.map((r) => (
                  <tr key={r.id}>
                    <td>{r.resource_code}</td>
                    <td>{r.resource_type}</td>
                    <td>{r.location || "—"}</td>
                    <td><StatusBadge
                      value={r.availability_status}
                      text={r.status_text}
                      tone={r.availability_status === "AVAILABLE" ? "ok" : r.availability_status === "OFFLINE" ? "muted" : "warn"}
                    /></td>
                    <td>{r.owner_team || "—"}</td>
                    <td>{r.maintenance_due_at ? formatDate(r.maintenance_due_at) : "—"}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>

      <section className="panel">
        <h2>预约日历 / 占用时间窗</h2>
        {bookings.length === 0 ? <EmptyState title="暂无预约" /> : <ResourceCalendar bookings={bookings} />}
      </section>
    </section>
  );
}
