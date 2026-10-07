export type GroundTaskStatus = "PENDING" | "ACCEPTED" | "IN_PROGRESS" | "BLOCKED" | "COMPLETED";

// 与后端一致：超时比较基准
export const CompareBasis = {
  Now: "NOW",
  ActualFinish: "ACTUAL_FINISH"
} as const;

export interface GroundTask {
  id: number;
  turnaround_id: number;
  flight_no: string;
  task_type: string;
  task_type_text: string;
  team_id: string;
  planned_start: string | null;
  /** 原计划截止时间（保留不变，是顺延基准） */
  deadline: string;
  /** 当前生效截止时间 = deadline + 未关闭延误分钟；超时一律按它判 */
  effective_deadline: string;
  open_delay_minutes: number;
  open_delay_count: number;
  accepted_at: string | null;
  actual_finish: string | null;
  status: GroundTaskStatus;
  status_text: string;
  blocker_note: string;
  overdue: boolean;
  overdue_minutes: number;
  compare_basis: string;
}
