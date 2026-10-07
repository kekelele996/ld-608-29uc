export const LOG_TEMPLATES = {
  FlightTurnaround: ["航班过站创建", "航班过站更新", "航班过站状态变更", "航班过站导出"],
  GroundTask: ["地勤任务创建", "地勤任务更新", "地勤任务状态变更", "地勤任务导出", "地勤任务签收"],
  GroundResource: ["保障资源创建", "保障资源更新", "保障资源状态变更", "保障资源导出"],
  ResourceBooking: ["资源预约创建", "资源预约更新", "资源预约状态变更", "资源预约导出"],
  DelayEvent: ["延误事件创建", "延误事件更新", "延误事件状态变更", "延误事件导出", "延误事件登记", "延误事件关闭"]
};

// 写操作日志模板索引，避免调用处散落魔法下标
export const TASK_LOG_SIGN_OFF = LOG_TEMPLATES.GroundTask[4];
export const DELAY_LOG_REGISTER = LOG_TEMPLATES.DelayEvent[4];
export const DELAY_LOG_CLOSE = LOG_TEMPLATES.DelayEvent[5];
