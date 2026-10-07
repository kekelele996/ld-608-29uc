import type { DelayEvent } from "../types/DelayEvent";
import type { GroundTask } from "../types/GroundTask";

/**
 * 延误顺延 / 超时判定的唯一前端口径，与后端 utils/delayPolicy.go 同构。
 * 看板超时卡片与任务列表都必须调用这里，禁止各页面各写一套。
 *
 * 口径（页面“判定口径”说明照此展示）：
 * 1. 任务 deadline（原计划截止时间）保留不变，登记延误不改写它；
 * 2. 另算当前生效截止时间 effective_deadline = deadline + 该航班「未关闭」延误分钟之和；
 *    关闭过的延误（resolved_at 非空）不参与顺延；
 * 3. 超时一律按 effective_deadline 判：
 *    - 已完成任务（actual_finish 有值）：actual_finish > effective_deadline；
 *    - 未完成任务（含仅签收、进行中、待签收）：now > effective_deadline。
 */

export const COMPARE_BASIS_NOW = "NOW";
export const COMPARE_BASIS_ACTUAL = "ACTUAL_FINISH";

export interface OpenDelaySummary {
  openMinutes: number;
  openCount: number;
}

/** 只统计尚未关闭（resolved_at 为空）的延误事件。 */
export function summarizeOpenDelays(events: DelayEvent[]): OpenDelaySummary {
  return events.reduce<OpenDelaySummary>(
    (acc, e) => {
      if (!e.resolved_at && !e.closed) {
        acc.openMinutes += e.minutes;
        acc.openCount += 1;
      }
      return acc;
    },
    { openMinutes: 0, openCount: 0 }
  );
}

/** 生效截止 = 原计划截止 + 未关闭延误分钟。 */
export function effectiveDeadline(deadlineIso: string, openMinutes: number, now: Date = new Date()): Date {
  const base = new Date(deadlineIso);
  // 无效时间兜底为当前，避免页面出现 Invalid Date 误判。
  if (Number.isNaN(base.getTime())) return now;
  return new Date(base.getTime() + openMinutes * 60_000);
}

export interface OverdueVerdict {
  overdue: boolean;
  overdueMinutes: number;
  compareBasis: string;
  effectiveAt: Date;
}

function diffMinutes(a: Date, b: Date): number {
  if (a.getTime() <= b.getTime()) return 0;
  return Math.round((a.getTime() - b.getTime()) / 60_000);
}

/**
 * 计算任务超时结论。
 * @param task 任务（deadline 原计划，open_delay_minutes 已含未关闭延误汇总）
 */
export function evaluateOverdue(task: GroundTask, now: Date = new Date()): OverdueVerdict {
  const effective = effectiveDeadline(task.deadline, task.open_delay_minutes, now);
  if (task.actual_finish) {
    const finished = new Date(task.actual_finish);
    return {
      overdue: finished.getTime() > effective.getTime(),
      overdueMinutes: diffMinutes(finished, effective),
      compareBasis: COMPARE_BASIS_ACTUAL,
      effectiveAt: effective
    };
  }
  return {
    overdue: now.getTime() > effective.getTime(),
    overdueMinutes: diffMinutes(now, effective),
    compareBasis: COMPARE_BASIS_NOW,
    effectiveAt: effective
  };
}

/** 给一批任务统一打上超时结论（看板卡片/列表共用）。 */
export function decorateTasks(tasks: GroundTask[], now: Date = new Date()) {
  return tasks.map((task) => {
    const v = evaluateOverdue(task, now);
    return {
      ...task,
      effective_deadline: v.effectiveAt.toISOString(),
      overdue: v.overdue,
      overdue_minutes: v.overdueMinutes,
      compare_basis: v.compareBasis
    } as GroundTask;
  });
}
