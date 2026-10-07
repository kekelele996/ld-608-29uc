// 与后端 constants/logTemplates.go 对齐的日志动作模板（前端用于按钮文案/操作日志展示）。
export const LOG_TEMPLATES = {
  TASK_DISPATCH: "派工：{taskType} → {team}，原计划截止 {deadline}",
  TASK_ACCEPT: "签收任务 #{id}，签收时间 {acceptedAt}，生效截止 {effectiveDeadline}",
  TASK_COMPLETE: "完成任务 #{id}，实际完成 {actualFinish}，超时 {overdue}",
  TASK_BLOCK: "阻塞任务 #{id}：{note}",
  DELAY_REGISTER: "登记延误 {delayType} {minutes} 分钟，未关闭累计顺延 {open} 分钟",
  DELAY_RESOLVE: "关闭延误 #{id}，{minutes} 分钟不再顺延",
  BOOKING_CONFLICT: "资源 {resource} 在该时间窗存在预约冲突"
} as const;
