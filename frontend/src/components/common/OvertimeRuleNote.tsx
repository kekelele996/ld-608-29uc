import { OVERTIME_RULE_LINES, OVERTIME_RULE_TITLE } from "../../constants/overtimeRule";

/** 超时判定口径说明条：看板与任务页共用，页面写清照哪条口径走。 */
export function OvertimeRuleNote() {
  return (
    <section className="rule-note">
      <strong>{OVERTIME_RULE_TITLE}</strong>
      <ol>
        {OVERTIME_RULE_LINES.map((line) => (
          <li key={line}>{line}</li>
        ))}
      </ol>
    </section>
  );
}
