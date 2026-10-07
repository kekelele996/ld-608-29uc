export function DelayTag({ minutes, closed = false }: { minutes: number; closed?: boolean }) {
  if (closed) return <span className="delay-tag closed">已关闭</span>;
  return <span className="delay-tag">顺延 +{minutes} 分钟</span>;
}
