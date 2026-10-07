import type { ResourceBooking } from "../../types/ResourceBooking";
import { ConflictBadge } from "./ConflictBadge";
import { StatusBadge } from "./StatusBadge";
import { formatClock } from "../../utils/formatters";

/** 资源占用日历（时间条）：按资源分组展示预约窗口与冲突。 */
export function ResourceCalendar({ bookings }: { bookings: ResourceBooking[] }) {
  const groups = bookings.reduce<Record<string, ResourceBooking[]>>((acc, b) => {
    (acc[b.resource_code] ??= []).push(b);
    return acc;
  }, {});
  return (
    <div className="calendar">
      {Object.entries(groups).map(([resource, list]) => (
        <div key={resource} className="calendar-row">
          <span className="calendar-res">{resource}</span>
          <div className="calendar-lane">
            {list.map((b) => (
              <span
                key={b.id}
                className={`calendar-bar ${b.booking_status === "CONFLICT" ? "conflict" : b.booking_status === "RELEASED" ? "released" : "held"}`}
                title={`${formatClock(b.start_time)} ~ ${formatClock(b.end_time)}`}
              >
                {b.booking_status === "CONFLICT" ? <ConflictBadge reason={b.conflict_reason} /> : <StatusBadge value={b.booking_status} text={b.status_text} />}
                <em>航班#{b.turnaround_id} {formatClock(b.start_time)}-{formatClock(b.end_time)}</em>
              </span>
            ))}
          </div>
        </div>
      ))}
    </div>
  );
}
