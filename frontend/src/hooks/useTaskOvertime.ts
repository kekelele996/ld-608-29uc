import { useMemo } from "react";
import { useGroundTaskStore } from "../stores/GroundTaskStore";
import { useDelayEventStore } from "../stores/DelayEventStore";
import { effectiveDeadline, isDelayOpen, openDelayMinutes, overtimeMinutes } from "../utils/overtime";
import type { GroundTask } from "../types/GroundTask";

export interface TaskOvertimeRow {
  task: GroundTask;
  /** 该任务所在航班当前未关闭延误累计的顺延分钟数 */
  delayMinutes: number;
  /** 当前生效截止时间（原计划截止 + 未关闭延误分钟） */
  effectiveDeadline: string;
  /** 超时分钟数，0 表示未超时 */
  overtimeMinutes: number;
  isOvertime: boolean;
}

/**
 * 看板超时卡片与任务列表共用的超时口径 hook：
 * 两个页面都从这里取数，保证超时条数与延误顺延分钟数一致。
 */
export function useTaskOvertime() {
  const tasks = useGroundTaskStore((state) => state.rows);
  const delays = useDelayEventStore((state) => state.rows);

  return useMemo(() => {
    const evaluatedAt = new Date();
    const rows: TaskOvertimeRow[] = tasks.map((task) => {
      const minutes = overtimeMinutes(task, delays, evaluatedAt);
      return {
        task,
        delayMinutes: openDelayMinutes(task.turnaround_id, delays),
        effectiveDeadline: effectiveDeadline(task, delays),
        overtimeMinutes: minutes,
        isOvertime: minutes > 0
      };
    });
    return {
      rows,
      evaluatedAt,
      overtimeRows: rows.filter((row) => row.isOvertime),
      overtimeCount: rows.filter((row) => row.isOvertime).length,
      openDelayCount: delays.filter(isDelayOpen).length,
      openDelayMinutesTotal: delays
        .filter(isDelayOpen)
        .reduce((sum, event) => sum + (Number(event.minutes) || 0), 0)
    };
  }, [tasks, delays]);
}
