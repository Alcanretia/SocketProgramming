import { ChatProvider, useChat } from "./socket/store";
import { ChatRoom } from "./components/shared/ChatRoom";
import { NotificationSetting } from "./components/shared/NotificationSetting";
import { roomInfo } from "./components/shared/roomInfo";
import { NameEntry } from "./features/users-dm/NameEntry";
import { UserList } from "./features/users-dm/UserList";
import { GroupPanel } from "./features/groups/GroupPanel";
import { InviteJoinDialog } from "./features/groups/InviteJoinDialog";
import { SummaryCard } from "./features/ai/SummaryCard";
import { UnreadPanel } from "./features/ai/UnreadPanel";

function Toasts() {
  const { state, actions } = useChat();
  return (
    <div className="toasts" aria-live="polite">
      {state.toasts.map((t) => (
        <div
          key={t.id}
          className={`toast ${t.kind}`}
          onClick={() => {
            if (t.room) actions.openRoom(t.room);
            actions.dismissToast(t.id);
          }}
        >
          {t.text}
        </div>
      ))}
    </div>
  );
}

function Shell() {
  const { state, actions } = useChat();

  if (!state.me) {
    return (
      <main className="center-page">
        {state.inviteToken && <p className="notice">You were invited to a group. Choose a name to continue.</p>}
        <NameEntry />
        {state.status !== "open" && <p className="muted">Connecting to the server...</p>}
        <Toasts />
      </main>
    );
  }

  const room = state.activeRoom;
  const info = room ? roomInfo(state, room) : null;
  const summary = room ? state.summaries[room] : undefined;

  return (
    <div className={`layout ${room ? "has-room" : ""}`}>
      {state.status !== "open" && <div className="banner">Connection lost. Reconnecting...</div>}
      <aside className="sidebar">
        <div className="me">
          Signed in as <strong>{state.me}</strong>
        </div>
        <UnreadPanel />
        <UserList />
        <GroupPanel />
        <NotificationSetting />
      </aside>
      <main className="main">
        {room && info ? (
          <ChatRoom
            title={info.title}
            subtitle={info.subtitle}
            me={state.me}
            messages={state.messages[room] ?? []}
            disabledReason={info.disabledReason}
            onSendText={(t) => actions.sendText(room, t)}
            onSendSticker={(id) => actions.sendSticker(room, id)}
            onBack={() => actions.openRoom(null)}
            header={summary && <SummaryCard state={summary} onDismiss={() => actions.dismissSummary(room)} />}
          />
        ) : (
          <div className="placeholder muted">Pick a user or a group to start chatting.</div>
        )}
      </main>
      <InviteJoinDialog />
      <Toasts />
    </div>
  );
}

export default function App() {
  return (
    <ChatProvider>
      <Shell />
    </ChatProvider>
  );
}
