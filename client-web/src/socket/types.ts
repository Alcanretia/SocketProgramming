// Types mirror docs/PROTOCOL.md. If this file changes, PROTOCOL.md changes first.

export type Kind = "text" | "sticker";

export interface ChatMessage {
  room: string;
  from: string;
  to?: string;
  groupId?: string;
  kind: Kind;
  text?: string;
  stickerId?: string;
  ts: number;
  msgId: string;
}

export interface GroupInfo {
  groupId: string;
  name: string;
  members: string[];
}

export interface ServerFrame {
  type: string;
  payload: Record<string, unknown>;
  reqId?: string;
}

export interface InviteInfo {
  groupId: string;
  groupName: string;
  memberCount: number;
  alreadyMember: boolean;
}

export type SummaryState =
  | { status: "loading" }
  | { status: "done"; text: string; messageCount: number }
  | { status: "error"; code: string; message: string };

/** Same formula as the server: lowercase both names, sort, join. */
export function dmRoomId(a: string, b: string): string {
  const x = a.toLowerCase();
  const y = b.toLowerCase();
  return x < y ? `dm:${x}|${y}` : `dm:${y}|${x}`;
}

export function groupRoomId(groupId: string): string {
  return `group:${groupId}`;
}

/** The other party's lowercase name in a DM room id, or null if not a DM. */
export function dmPeerKey(room: string, me: string): string | null {
  if (!room.startsWith("dm:")) return null;
  const [a, b] = room.slice(3).split("|");
  const mine = me.toLowerCase();
  return a === mine ? b : b === mine ? a : null;
}
