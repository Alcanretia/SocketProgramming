import type { State } from "../../socket/store";
import { dmPeerKey } from "../../socket/types";

export interface RoomInfo {
  title: string;
  subtitle?: string;
  disabledReason?: string;
}

/** Display data for a room id. UI hints only: the server enforces who may post or read. */
export function roomInfo(s: State, room: string): RoomInfo {
  const me = s.me ?? "";
  const peer = dmPeerKey(room, me);
  if (peer) {
    const online = s.users.find((u) => u.toLowerCase() === peer);
    return {
      title: online ?? peer,
      subtitle: "Direct message",
      disabledReason: online ? undefined : "This user is offline.",
    };
  }
  if (room.startsWith("group:")) {
    const g = s.groups.find((x) => x.groupId === room.slice(6));
    if (!g) return { title: "Group", disabledReason: "This group no longer exists." };
    const member = g.members.some((m) => m.toLowerCase() === me.toLowerCase());
    return {
      title: g.name,
      subtitle: g.members.join(", "),
      disabledReason: member ? undefined : "Join this group to chat.",
    };
  }
  return { title: room };
}
