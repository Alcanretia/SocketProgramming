import { useEffect, useState } from "react";
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

const TOAST_LIFETIME_MS = 6000; // how long the oldest toast stays before it goes away
const TOAST_FADE_MS = 400; // the last part of that time is spent fading out (keep in sync with styles.css)

function Toasts() {
  const { state, actions } = useChat();
  const { toasts } = state;
  const [fadingId, setFadingId] = useState<number | null>(null);

  // One timer, always for the OLDEST toast: when it expires that toast fades and
  // is removed, and then the next-oldest gets its own full countdown. A newly
  // arrived toast changes `newestId`, which restarts the countdown of the oldest
  // one, so a busy chat keeps its toasts on screen. (Ids only grow, so the newest
  // id changes exactly when something arrives, or when the newest is clicked away.)
  const oldestId = toasts.length ? toasts[0].id : null;
  const newestId = toasts.length ? toasts[toasts.length - 1].id : null;
  const { dismissToast } = actions;
  useEffect(() => {
    setFadingId(null);
    if (oldestId === null) return;
    const fade = window.setTimeout(() => setFadingId(oldestId), TOAST_LIFETIME_MS - TOAST_FADE_MS);
    const remove = window.setTimeout(() => dismissToast(oldestId), TOAST_LIFETIME_MS);
    return () => {
      window.clearTimeout(fade);
      window.clearTimeout(remove);
    };
  }, [oldestId, newestId, dismissToast]);

  return (
    <div className="toasts" aria-live="polite">
      {toasts.map((t) => (
        <div
          key={t.id}
          className={`toast ${t.kind}${t.id === fadingId ? " fading" : ""}`}
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
