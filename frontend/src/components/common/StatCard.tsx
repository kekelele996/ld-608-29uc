import type { ReactNode } from "react";

export function StatCard({ label, value, tone, hint }: { label: string; value: ReactNode; tone?: "danger" | "warn" | "ok"; hint?: string }) {
  return (
    <div className={`stat ${tone ? `stat-${tone}` : ""}`}>
      <span>{label}</span>
      <strong>{value}</strong>
      {hint ? <em>{hint}</em> : null}
    </div>
  );
}
