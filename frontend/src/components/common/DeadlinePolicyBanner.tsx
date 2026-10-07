/**
 * 超时判定口径说明 —— 需求要求“页面写清照哪条走”。
 * 看板、任务页都渲染它，避免两套说法。
 */
export function DeadlinePolicyBanner() {
  return (
    <div className="policy-banner">
      <strong>超时判定口径</strong>
      <ol>
        <li>任务 <b>原计划截止时间（deadline）保留不变</b>，登记/关闭延误都不改写该列；</li>
        <li>另算<b>当前生效截止时间</b> = 原计划截止 + 该航班<b>尚未关闭</b>延误事件分钟之和；关闭过的延误不参与顺延；</li>
        <li>超时一律按<b>生效截止</b>判：未签收/进行中任务与「当前时间」比；已签收完成的任务与「实际完成时间」比。</li>
      </ol>
      <p className="policy-note">看板“超时任务”卡片与任务列表共用同一计算结果（utils/delayPolicy ↔ 后端 ViewService）。</p>
    </div>
  );
}
