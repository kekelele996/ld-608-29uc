import type { FlightTurnaround } from "../../types/FlightTurnaround";
import { StatusBadge } from "./StatusBadge";
import { DelayTag } from "./DelayTag";
import { formatClock } from "../../utils/formatters";

const STAGES = ["ARRIVING", "ON_STAND", "IN_SERVICE", "READY", "DEPARTED"] as const;

/** 航班过站时间轴：展示状态推进与未关闭延误累计。 */
export function TurnaroundTimeline({ flight }: { flight: FlightTurnaround }) {
  const activeIndex = flight.turnaround_status === "DELAYED"
    ? STAGES.indexOf("IN_SERVICE")
    : STAGES.indexOf(flight.turnaround_status as (typeof STAGES)[number]);
  return (
    <div className="timeline">
      <div className="timeline-head">
        <strong>{flight.flight_no} · {flight.stand_no}</strong>
        <StatusBadge value={flight.turnaround_status} text={flight.status_text} tone={flight.turnaround_status === "DELAYED" ? "warn" : "ok"} />
      </div>
      <ol className="timeline-track">
        {STAGES.map((stage, i) => (
          <li key={stage} className={i <= activeIndex ? "done" : "todo"}>
            <i />
            <span>{stage.replace(/_/g, " ")}</span>
          </li>
        ))}
      </ol>
      <footer>
        <span>到港 {formatClock(flight.arrival_time)} · 计划离港 {formatClock(flight.departure_time)}</span>
        {flight.open_delay_minutes > 0 ? <DelayTag minutes={flight.open_delay_minutes} label="未关闭顺延" /> : null}
      </footer>
    </div>
  );
}
