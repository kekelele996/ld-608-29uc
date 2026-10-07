import { ERROR_CODES, type ErrorCode } from "./errorCodes";

export const ERROR_MESSAGES: Record<ErrorCode, string> = {
  AUTH_REQUIRED: "请先登录后再继续操作",
  AUTH_INVALID: "用户名或口令不正确（演示账号 demo123）",
  FORBIDDEN: "当前角色没有执行该动作的权限",
  RBAC_DENIED: "当前角色没有执行该动作的权限",
  VALIDATION_FAILED: "表单字段缺失或格式错误",
  RECORD_NOT_FOUND: "目标记录不存在",
  INVALID_ENUM: "枚举值不合法",
  RESOURCE_CONFLICT: "资源在该时间窗已被预约，存在冲突",
  DEADLINE_NOT_EXTENDED: "延误分钟必须为正数才会顺延截止时间",
  TASK_NOT_ACCEPTABLE: "任务当前状态不允许该操作",
  RATE_LIMITED: "请求过于频繁，请稍后再试",
  NETWORK_OFFLINE: "后端不可达，已切换本地演示数据",
  INTERNAL_ERROR: "服务内部错误"
};

export { ERROR_CODES };
