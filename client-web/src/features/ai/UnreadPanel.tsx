import { useChat } from "../../socket/store";
import { roomInfo } from "../../components/shared/roomInfo";
import { SummaryCard } from "./SummaryCard";

/**
 * Rooms with unread messages (counts come from the server) and the AI summary
 * button (special feature F1). The summary request goes over the WebSocket; the
 * server builds the prompt only from messages this user was allowed to receive
 * and calls the LLM, so the API key never reaches the browser.
 */
export function UnreadPanel() {
  const { state, actions } = useChat();
  const rooms = Array.from(
    new Set([
      ...Object.entries(state.unread)
        .filter(([, n]) => n > 0)
        .map(([r]) => r),
      ...Object.keys(state.summaries).filter((r) => r !== state.activeRoom),
    ]),
  );
  if (rooms.length === 0) return null;

  return (
    <section>
      <h2>Unread</h2>
      <ul className="list">
        {rooms.map((room) => {
          const n = state.unread[room] ?? 0;
          const sum = state.summaries[room];
          return (
            <li key={room} className="unread-item">
              <div className="row static">
                <span className="grow">{roomInfo(state, room).title}</span>
                {n > 0 && <span className="badge">{n}</span>}
              </div>
              <div className="actions">
                <button className="small secondary" onClick={() => actions.openRoom(room)}>
                  Open
                </button>
                {n > 0 && (
                  <button className="small" disabled={sum?.status === "loading"} onClick={() => actions.summarize(room)}>
                    ✨ Summarize
                  </button>
                )}
              </div>
              {sum && <SummaryCard state={sum} onDismiss={() => actions.dismissSummary(room)} />}
            </li>
          );
        })}
      </ul>
      <p className="muted small">Summaries send the unread messages of that chat to an AI service (Anthropic).</p>
    </section>
  );
}
