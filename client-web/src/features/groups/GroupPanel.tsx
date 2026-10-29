import { useState, type FormEvent } from "react";
import { useChat } from "../../socket/store";
import { groupRoomId } from "../../socket/types";
import { InviteQr } from "./InviteQr";

/** Create a group (R8), list every group with its members (R9), join by yourself (R10). */
export function GroupPanel() {
  const { state, actions } = useChat();
  const [name, setName] = useState("");
  const [qrFor, setQrFor] = useState<string | null>(null);
  const me = (state.me ?? "").toLowerCase();

  const create = (e: FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;
    actions.createGroup(name.trim());
    setName("");
  };

  return (
    <section>
      <h2>Groups ({state.groups.length})</h2>
      <form className="inline-form" onSubmit={create}>
        <input value={name} maxLength={30} onChange={(e) => setName(e.target.value)} placeholder="New group name" aria-label="New group name" />
        <button type="submit" disabled={!name.trim()}>
          Create
        </button>
      </form>

      <ul className="list">
        {state.groups.map((g) => {
          const member = g.members.some((m) => m.toLowerCase() === me);
          const room = groupRoomId(g.groupId);
          const unread = state.unread[room] ?? 0;
          return (
            <li key={g.groupId} className="group-item">
              <button className={`row ${state.activeRoom === room ? "active" : ""}`} disabled={!member} onClick={() => actions.openRoom(room)}>
                <span className="grow">
                  <strong>{g.name}</strong>
                  <span className="muted small members">{g.members.join(", ")}</span>
                </span>
                {unread > 0 && <span className="badge">{unread}</span>}
              </button>
              <div className="actions">
                {member ? (
                  <button className="small" onClick={() => setQrFor(qrFor === g.groupId ? null : g.groupId)}>
                    {qrFor === g.groupId ? "Hide QR" : "Invite (QR)"}
                  </button>
                ) : (
                  <button className="small" onClick={() => actions.joinGroup(g.groupId)}>
                    Join
                  </button>
                )}
              </div>
              {qrFor === g.groupId && member && <InviteQr groupId={g.groupId} />}
            </li>
          );
        })}
        {state.groups.length === 0 && <li className="muted small">No groups yet.</li>}
      </ul>
    </section>
  );
}
