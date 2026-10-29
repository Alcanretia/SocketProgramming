import type { SummaryState } from "../../socket/types";

const REASONS: Record<string, string> = {
  nothing_unread: "There is nothing unread to summarize.",
  timeout: "The AI service took too long. Try again.",
  llm_error: "The AI summary is not available right now.",
  not_allowed: "You cannot summarize this chat.",
};

export function SummaryCard({ state, onDismiss }: { state: SummaryState; onDismiss: () => void }) {
  return (
    <div className="summary-card" role="status">
      <div className="summary-head">
        <strong>✨ AI summary</strong>
        <button className="small secondary" onClick={onDismiss} aria-label="Dismiss summary">
          ×
        </button>
      </div>
      {state.status === "loading" && <p className="muted">Summarizing...</p>}
      {state.status === "error" && <p className="error">{REASONS[state.code] ?? state.message}</p>}
      {state.status === "done" && (
        <>
          {/* pre-wrap keeps the model's bullet lines; React still escapes the text */}
          <p className="summary-text">{state.text}</p>
          <p className="muted small">Based on {state.messageCount} unread messages.</p>
        </>
      )}
    </div>
  );
}
