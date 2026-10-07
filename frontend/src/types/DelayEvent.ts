export interface DelayEvent {
  id: number;
  turnaround_id: number;
  flight_no: string;
  delay_type: string;
  minutes: number;
  root_cause: string;
  responsibility_team: string;
  created_at: string;
  /** 关闭时间；null 表示尚未关闭，仍在顺延任务截止 */
  resolved_at: string | null;
  closed: boolean;
  /** 该航班当前未关闭延误累计分钟 */
  open_delay_minutes: number;
}

export interface DelayRegisterPayload {
  turnaround_id: number;
  delay_type: string;
  minutes: number;
  root_cause?: string;
  responsibility_team?: string;
}
