import { useChat } from "../../socket/store";
import { dmRoomId } from "../../socket/types";

/** Every connected name, including our own (R4). Click another name to open a DM room (R5). */
export function UserList() {
  const { state, actions } = useChat();
  const me = state.me ?? "";

  return (
    <section>
      <h2>Online ({state.users.length})</h2>
      <ul className="list">
        {state.users.map((u) => {
          const isMe = u.toLowerCase() === me.toLowerCase();
          const room = dmRoomId(me, u);
          const unread = state.unread[room] ?? 0;
          return (
            <li key={u}>
              <button
                className={`row ${state.activeRoom === room ? "active" : ""}`}
                disabled={isMe}
                onClick={() => actions.openRoom(room)}
              >
                <span className="dot" />
                <span className="grow">
                  {u}
                  {isMe && <em className="muted"> (you)</em>}
                </span>
                {unread > 0 && <span className="badge">{unread}</span>}
              </button>
            </li>
          );
        })}
      </ul>
    </section>
  );
}
