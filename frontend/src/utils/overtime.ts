import type { GroundTask } from "../types/GroundTask";
import type { DelayEvent } from "../types/DelayEvent";

// 超时判定口径（与后端 services/GroundTask.go 镜像保持一致）：
// 1. 任务原计划截止时间 deadline 不因登记/关闭延误而改写；
// 2. 当前生效截止时间 = deadline + 该航班「尚未关闭」的延误分钟数之和；
// 3. 已签收/已完成的任务（actual_finish 非空）按实际完成时间比较，其余按当前时间比较。

const MINUTE_MS = 60_000;

const parseTime = (value: string): number => {
  if (!value) return Number.NaN;
  return Date.parse(value);
};

/** 延误事件是否仍未关闭：resolved_at 为空表示未关闭，未关闭才参与顺延。 */
export const isDelayOpen = (event: DelayEvent): boolean =>
  Number.isNaN(parseTime(event.resolved_at));

/** 某航班当前仍生效的顺延分钟数：仅累加未关闭延误。 */
export function openDelayMinutes(turnaroundId: number, delays: DelayEvent[]): number {
  return delays
    .filter((event) => event.turnaround_id === turnaroundId && isDelayOpen(event))
    .reduce((sum, event) => sum + (Number(event.minutes) || 0), 0);
}

/** 当前生效截止时间（ISO 字符串）：原计划截止 + 未关闭延误分钟。 */
export function effectiveDeadline(task: GroundTask, delays: DelayEvent[]): string {
  const base = parseTime(task.deadline);
  if (Number.isNaN(base)) return task.deadline;
  const shifted = base + openDelayMinutes(task.turnaround_id, delays) * MINUTE_MS;
  return new Date(shifted).toISOString();
}

/** 任务是否已有实际完成时间（签收或完成都会写入 actual_finish）。 */
export const hasActualFinish = (task: GroundTask): boolean =>
  !Number.isNaN(parseTime(task.actual_finish));

/**
 * 超时分钟数，0 表示未超时。
 * 已签收/已完成：actual_finish 与生效截止比较；未完成：评估时刻与生效截止比较。
 */
export function overtimeMinutes(task: GroundTask, delays: DelayEvent[], now: Date = new Date()): number {
  const effective = parseTime(effectiveDeadline(task, delays));
  if (Number.isNaN(effective)) return 0;
  const comparedAt = hasActualFinish(task) ? parseTime(task.actual_finish) : now.getTime();
  if (Number.isNaN(comparedAt)) return 0;
  return Math.max(0, Math.round((comparedAt - effective) / MINUTE_MS));
}

export function isTaskOvertime(task: GroundTask, delays: DelayEvent[], now: Date = new Date()): boolean {
  return overtimeMinutes(task, delays, now) > 0;
}
