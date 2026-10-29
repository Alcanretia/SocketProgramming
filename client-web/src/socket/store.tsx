import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useReducer,
  useRef,
  useState,
  type ReactNode,
} from "react";
import { Connection, type ConnStatus } from "./connection";
import { beep, setTitleBadge, showSystemNotification } from "./notifications";
import {
  dmPeerKey,
  type ChatMessage,
  type GroupInfo,
  type InviteInfo,
  type ServerFrame,
  type SummaryState,
} from "./types";

export interface Toast {
  id: number;
  kind: "info" | "error" | "message";
  text: string;
  room?: string;
}

export interface State {
  status: ConnStatus;
  me: string | null;
  wantedName: string | null; // name we (re-)register as whenever the socket is open
  registerError: string | null;
  users: string[];
  groups: GroupInfo[];
  messages: Record<string, ChatMessage[]>;
  unread: Record<string, number>; // owned by the server, mirrored here
  activeRoom: string | null;
  toasts: Toast[];
  summaries: Record<string, SummaryState>;
  inviteToken: string | null; // from ?invite=<token> in the URL
  invitePreview: InviteInfo | null;
  inviteCreated: Record<string, string>; // groupId -> token
}

type Action =
  | { t: "status"; status: ConnStatus }
  | { t: "frame"; frame: ServerFrame }
  | { t: "wantName"; name: string }
  | { t: "setActive"; room: string | null }
  | { t: "toast"; toast: Toast }
  | { t: "dismissToast"; id: number }
  | { t: "summaryLoading"; room: string }
  | { t: "summaryDismiss"; room: string }
  | { t: "inviteClear" };

let toastSeq = 0;
const mkToast = (kind: Toast["kind"], text: string, room?: string): Toast => ({
  id: ++toastSeq,
  kind,
  text,
  room,
});

function initState(): State {
  let token: string | null = null;
  try {
    token = new URLSearchParams(location.search).get("invite");
  } catch {
    token = null;
  }
  return {
    status: "connecting",
    me: null,
    wantedName: null,
    registerError: null,
    users: [],
    groups: [],
    messages: {},
    unread: {},
    activeRoom: null,
    toasts: [],
    summaries: {},
    inviteToken: token,
    invitePreview: null,
    inviteCreated: {},
  };
}

const payload = <T,>(f: ServerFrame) => f.payload as unknown as T;

function reducer(s: State, a: Action): State {
  switch (a.t) {
    case "status":
      if (a.status === "open") return { ...s, status: "open" };
      // The server forgets everything about us when the socket closes.
      return { ...s, status: a.status, me: null, users: [], groups: [], unread: {} };
    case "wantName":
      return { ...s, wantedName: a.name, registerError: null };
    case "setActive":
      return { ...s, activeRoom: a.room };
    case "toast":
      return { ...s, toasts: [...s.toasts.slice(-3), a.toast] };
    case "dismissToast":
      return { ...s, toasts: s.toasts.filter((t) => t.id !== a.id) };
    case "summaryLoading":
      return { ...s, summaries: { ...s.summaries, [a.room]: { status: "loading" } } };
    case "summaryDismiss": {
      const { [a.room]: _gone, ...rest } = s.summaries;
      return { ...s, summaries: rest };
    }
    case "inviteClear":
      return { ...s, inviteToken: null, invitePreview: null };
    case "frame":
      return onFrame(s, a.frame);
  }
}

function onFrame(s: State, f: ServerFrame): State {
  switch (f.type) {
    case "register_ok":
      return { ...s, me: payload<{ name: string }>(f).name, registerError: null };
    case "register_error":
      // wantedName is cleared so we do not retry the same rejected name forever
      return { ...s, wantedName: null, registerError: payload<{ message: string }>(f).message };
    case "user_list":
      return { ...s, users: payload<{ users: string[] }>(f).users };
    case "group_list":
      return { ...s, groups: payload<{ groups: GroupInfo[] }>(f).groups };
    case "unread":
      return { ...s, unread: payload<{ counts: Record<string, number> }>(f).counts };
    case "dm":
    case "group_msg": {
      const m = payload<ChatMessage>(f);
      const list = s.messages[m.room] ?? [];
      if (list.some((x) => x.msgId === m.msgId)) return s;
      return { ...s, messages: { ...s.messages, [m.room]: [...list, m] } };
    }
    case "create_group_ok":
      return { ...s, activeRoom: `group:${payload<{ groupId: string }>(f).groupId}` };
    case "join_group_ok":
      return {
        ...s,
        activeRoom: `group:${payload<{ groupId: string }>(f).groupId}`,
        inviteToken: null,
        invitePreview: null,
      };
    case "invite_created": {
      const p = payload<{ groupId: string; token: string }>(f);
      return { ...s, inviteCreated: { ...s.inviteCreated, [p.groupId]: p.token } };
    }
    case "invite_info_ok":
      return { ...s, invitePreview: payload<InviteInfo>(f) };
    case "summary_result": {
      const p = payload<{ room: string; text: string; messageCount: number }>(f);
      return {
        ...s,
        summaries: { ...s.summaries, [p.room]: { status: "done", text: p.text, messageCount: p.messageCount } },
      };
    }
    case "summary_error": {
      const p = payload<{ room: string; code: string; message: string }>(f);
      return {
        ...s,
        summaries: { ...s.summaries, [p.room]: { status: "error", code: p.code, message: p.message } },
      };
    }
    case "error": {
      const p = payload<{ code: string; message: string }>(f);
      const badInvite = p.code === "bad_invite";
      return {
        ...s,
        toasts: [...s.toasts.slice(-3), mkToast("error", p.message)],
        inviteToken: badInvite ? null : s.inviteToken,
        invitePreview: badInvite ? null : s.invitePreview,
      };
    }
    default:
      return s;
  }
}

export interface Actions {
  register(name: string): void;
  openRoom(room: string | null): void;
  sendText(room: string, text: string): void;
  sendSticker(room: string, stickerId: string): void;
  createGroup(name: string): void;
  joinGroup(groupId: string, invite?: string): void;
  createInvite(groupId: string): void;
  summarize(room: string): void;
  dismissSummary(room: string): void;
  dismissToast(id: number): void;
  declineInvite(): void;
  notify(text: string): void;
}

interface Ctx {
  state: State;
  actions: Actions;
}

const ChatContext = createContext<Ctx | null>(null);

export function useChat(): Ctx {
  const ctx = useContext(ChatContext);
  if (!ctx) throw new Error("useChat must be used inside <ChatProvider>");
  return ctx;
}

export function ChatProvider({ children }: { children: ReactNode }) {
  const [state, dispatch] = useReducer(reducer, undefined, initState);
  const stateRef = useRef(state);
  stateRef.current = state;
  const conn = useRef<Connection | null>(null);
  const [visible, setVisible] = useState(!document.hidden);

  // One WebSocket for the whole app.
  useEffect(() => {
    const c = new Connection(Connection.defaultUrl(), {
      onStatus: (status) => dispatch({ t: "status", status }),
      onFrame: (frame) => {
        alertIfIncoming(frame, stateRef.current, (room) => dispatch({ t: "setActive", room }), (toast) =>
          dispatch({ t: "toast", toast }),
        );
        dispatch({ t: "frame", frame });
      },
    });
    conn.current = c;
    c.connect();
    return () => {
      c.close();
      conn.current = null;
    };
  }, []);

  // Track tab visibility: a hidden tab does not count as "viewing the room".
  useEffect(() => {
    const onVis = () => setVisible(!document.hidden);
    document.addEventListener("visibilitychange", onVis);
    return () => document.removeEventListener("visibilitychange", onVis);
  }, []);

  const send = useCallback((type: string, p: Record<string, unknown> = {}) => {
    conn.current?.send(type, p);
  }, []);

  // (Re-)register whenever the socket is open and we want a name but have none.
  useEffect(() => {
    if (state.status === "open" && state.wantedName && !state.me) {
      send("register", { name: state.wantedName });
    }
  }, [state.status, state.wantedName, state.me, send]);

  // Tell the server which room we are looking at so it counts unread correctly.
  useEffect(() => {
    const room = state.activeRoom;
    if (!state.me || !room || !visible || state.status !== "open") return;
    send("open_room", { room });
    return () => send("close_room", { room });
  }, [state.me, state.activeRoom, state.status, visible, send]);

  // Arrived via a QR code: ask the server what the invite is (this does NOT join).
  const askedInvite = useRef<string | null>(null);
  useEffect(() => {
    if (state.me && state.inviteToken && askedInvite.current !== state.inviteToken) {
      askedInvite.current = state.inviteToken;
      send("invite_info", { token: state.inviteToken });
    }
  }, [state.me, state.inviteToken, send]);

  // Remove ?invite= from the address bar once it has been handled.
  useEffect(() => {
    if (!state.inviteToken && location.search.includes("invite=")) {
      history.replaceState(null, "", location.pathname);
    }
  }, [state.inviteToken]);

  // Unread badge in the tab title.
  useEffect(() => {
    setTitleBadge(Object.values(state.unread).reduce((a, b) => a + b, 0));
  }, [state.unread]);

  const actions = useMemo<Actions>(() => {
    const chat = (room: string, body: Record<string, unknown>) => {
      const peer = dmPeerKey(room, stateRef.current.me ?? "");
      if (peer) send("dm", { to: peer, ...body });
      else if (room.startsWith("group:")) send("group_msg", { groupId: room.slice(6), ...body });
    };
    return {
      register: (name) => dispatch({ t: "wantName", name }),
      openRoom: (room) => dispatch({ t: "setActive", room }),
      sendText: (room, text) => chat(room, { kind: "text", text }),
      sendSticker: (room, stickerId) => chat(room, { kind: "sticker", stickerId }),
      createGroup: (name) => send("create_group", { name }),
      joinGroup: (groupId, invite) => send("join_group", invite ? { groupId, invite } : { groupId }),
      createInvite: (groupId) => send("create_invite", { groupId }),
      summarize: (room) => {
        dispatch({ t: "summaryLoading", room });
        send("summarize_request", { room });
      },
      dismissSummary: (room) => dispatch({ t: "summaryDismiss", room }),
      dismissToast: (id) => dispatch({ t: "dismissToast", id }),
      declineInvite: () => dispatch({ t: "inviteClear" }),
      notify: (text) => dispatch({ t: "toast", toast: mkToast("info", text) }),
    };
  }, [send]);

  return <ChatContext.Provider value={{ state, actions }}>{children}</ChatContext.Provider>;
}

// ---- incoming-message alerts (special feature F4) --------------------------

function alertIfIncoming(
  f: ServerFrame,
  s: State,
  openRoom: (room: string) => void,
  toast: (t: Toast) => void,
): void {
  if (f.type !== "dm" && f.type !== "group_msg") return;
  const m = f.payload as unknown as ChatMessage;
  if (!s.me || m.from === s.me) return; // never alert for our own messages

  const watching = s.activeRoom === m.room && !document.hidden;
  if (watching) return;

  const where =
    f.type === "group_msg" ? (s.groups.find((g) => g.groupId === m.groupId)?.name ?? "a group") : null;
  const title = where ? `${m.from} in ${where}` : m.from;
  const body = m.kind === "sticker" ? "sent a sticker" : (m.text ?? "");

  // 1) in-app toast: always works, even on a plain-HTTP LAN demo
  toast(mkToast("message", `${title}: ${body}`, m.room));
  beep();
  // 2) OS-level notification: only when the tab is hidden, and only if the
  //    browser allows it (secure context + permission)
  if (document.hidden) showSystemNotification(title, body, () => openRoom(m.room));
}
