import { useChat } from "../../socket/store";
import { dmRoomId } from "../../socket/types";

// ============================================================================
// P2 Frontend — งานย่อย 6.2: แถบรายชื่อผู้ใช้ออนไลน์ (UserList.tsx)
// ครอบคลุมเกณฑ์:
//   - R4: แสดงรายชื่อผู้ใช้ออนไลน์ทุกคน รวมถึงชื่อของตัวเอง (มีป้าย "(you)")
//   - R5: เมื่อคลิกที่ชื่อคนอื่น จะคำนวณรหัสห้อง DM ด้วย dmRoomId(me, u) แล้วเปิดห้องแชทนั้น
//         ส่วนปุ่มชื่อของตัวเองจะถูกล็อก (disabled={isMe}) ไม่ให้กดคุยกับตัวเองได้
// ============================================================================

export function UserList() {
  const { state, actions } = useChat();
  const me = state.me ?? "";

  return (
    <section className="panel" aria-label="Online users">
      <h2>Online ({state.users.length})</h2>
      <ul className="list">
        {state.users.map((u) => {
          const isMe = u.toLowerCase() === me.toLowerCase();
          const room = isMe ? null : dmRoomId(me, u);
          const unread = room ? (state.unread[room] ?? 0) : 0;
          const active = room !== null && state.activeRoom === room;

          return (
            <li key={u.toLowerCase()}>
              <button
                type="button"
                className={`list-item ${active ? "active" : ""}`}
                disabled={isMe}
                onClick={() => room && actions.openRoom(room)}
              >
                <span>
                  {u} {isMe && <em className="muted">(you)</em>}
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
