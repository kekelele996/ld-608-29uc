import { useEffect, useMemo, useState } from "react";
import { useDelayEventStore } from "../stores/DelayEventStore";
import { useFlightTurnaroundStore } from "../stores/FlightTurnaroundStore";
import { createDelayEventForm } from "../constructors/DelayEventConstructor";
import { isDelayOpen } from "../utils/overtime";
import { formatDateOrDash, formatMinutes } from "../utils/formatters";
import { GroundTaskType } from "../constants/GroundTaskType";
import { STATUS_TEXT } from "../constants/statusText";
import { DelayTag } from "../components/common/DelayTag";
import { EmptyState } from "../components/common/EmptyState";

const delayTypeText = (value: string) =>
  (STATUS_TEXT.GroundTaskType as Record<string, string>)[value] ?? value;

export function DelaysPage() {
  const delays = useDelayEventStore((state) => state.rows);
  const loadDelays = useDelayEventStore((state) => state.load);
  const register = useDelayEventStore((state) => state.register);
  const close = useDelayEventStore((state) => state.close);
  const turnarounds = useFlightTurnaroundStore((state) => state.rows);
  const loadTurnarounds = useFlightTurnaroundStore((state) => state.load);

  const [form, setForm] = useState(() => createDelayEventForm());

  useEffect(() => {
    void loadDelays();
    void loadTurnarounds();
  }, [loadDelays, loadTurnarounds]);

  const flightNo = useMemo(() => {
    const map = new Map(turnarounds.map((item) => [item.id, item.flight_no]));
    return (id: number) => map.get(id) ?? `#${id}`;
  }, [turnarounds]);

  const submit = async () => {
    if (!form.turnaround_id || form.minutes <= 0) return;
    await register(form);
    setForm(createDelayEventForm());
  };

  return (
    <main className="page">
      <section className="page-head">
        <div>
          <p className="eyebrow">ground-turn</p>
          <h1>延误归因</h1>
        </div>
      </section>

      <section className="panel wide">
        <h2>登记延误（登记后该航班任务的生效截止时间立即顺延）</h2>
        <div className="form-row">
          <label>
            航班
            <select
              value={form.turnaround_id}
              onChange={(event) => setForm({ ...form, turnaround_id: Number(event.target.value) })}
            >
              <option value={0}>请选择航班</option>
              {turnarounds.map((item) => (
                <option key={item.id} value={item.id}>
                  {item.flight_no}（机位 {item.stand_no}）
                </option>
              ))}
            </select>
          </label>
          <label>
            延误类型
            <select
              value={form.delay_type}
              onChange={(event) => setForm({ ...form, delay_type: event.target.value })}
            >
              {GroundTaskType.map((value) => (
                <option key={value} value={value}>
                  {delayTypeText(value)}
                </option>
              ))}
            </select>
          </label>
          <label>
            影响分钟数
            <input
              type="number"
              min={1}
              value={form.minutes || ""}
              onChange={(event) => setForm({ ...form, minutes: Number(event.target.value) })}
            />
          </label>
          <label>
            原因
            <input
              value={form.root_cause}
              onChange={(event) => setForm({ ...form, root_cause: event.target.value })}
            />
          </label>
          <label>
            责任班组
            <input
              value={form.responsibility_team}
              onChange={(event) => setForm({ ...form, responsibility_team: event.target.value })}
            />
          </label>
          <button className="btn primary" onClick={() => void submit()}>
            登记延误
          </button>
        </div>
      </section>

      <section className="panel wide">
        <h2>延误事件（关闭后不再参与任务截止顺延）</h2>
        {delays.length === 0 ? (
          <EmptyState title="暂无延误事件" />
        ) : (
          <table className="data-table">
            <thead>
              <tr>
                <th>航班</th>
                <th>类型</th>
                <th>影响分钟</th>
                <th>原因</th>
                <th>责任班组</th>
                <th>状态</th>
                <th>关闭时间</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              {delays.map((event) => {
                const open = isDelayOpen(event);
                return (
                  <tr key={event.id}>
                    <td>{flightNo(event.turnaround_id)}</td>
                    <td>{delayTypeText(event.delay_type)}</td>
                    <td>{formatMinutes(Number(event.minutes) || 0)}</td>
                    <td>{event.root_cause}</td>
                    <td>{event.responsibility_team}</td>
                    <td>
                      <DelayTag minutes={Number(event.minutes) || 0} closed={!open} />
                    </td>
                    <td>{formatDateOrDash(event.resolved_at)}</td>
                    <td>
                      {open ? (
                        <button className="btn" onClick={() => void close(event.id)}>
                          关闭
                        </button>
                      ) : (
                        <span className="muted">已关闭</span>
                      )}
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        )}
      </section>
    </main>
  );
}
