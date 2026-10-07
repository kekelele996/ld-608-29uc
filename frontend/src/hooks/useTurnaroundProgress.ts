import { useMemo } from "react";
import type { FlightTurnaround } from "../types/FlightTurnaround";

/** 单个/一批航班的任务完成率、超时率、未关闭延误汇总。看板与航班页共用。 */
export function useTurnaroundProgress(flights: FlightTurnaround[]) {
  return useMemo(() => {
    const totals = flights.reduce(
      (acc, f) => {
        acc.tasks += f.task_total;
        acc.completed += f.task_completed;
        acc.overdue += f.task_overdue;
        acc.openDelayMinutes += f.open_delay_minutes;
        acc.openDelayEvents += f.open_delay_count;
        return acc;
      },
      { tasks: 0, completed: 0, overdue: 0, openDelayMinutes: 0, openDelayEvents: 0 }
    );
    const completionRate = totals.tasks === 0 ? 0 : Math.round((totals.completed / totals.tasks) * 100);
    const overdueRate = totals.tasks === 0 ? 0 : Math.round((totals.overdue / totals.tasks) * 100);
    return { ...totals, completionRate, overdueRate };
  }, [flights]);
}
