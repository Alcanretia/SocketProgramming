import { useChat } from "../../socket/store";

/**
 * Shown after someone opens an invite link. Opening the link does not add
 * anyone: the user has to confirm here, and only then is `join_group` sent (R10).
 */
export function InviteJoinDialog() {
  const { state, actions } = useChat();
  const info = state.invitePreview;
  if (!state.inviteToken || !info) return null;

  const onPrimary = () => {
    if (info.alreadyMember) {
      actions.openRoom(`group:${info.groupId}`);
      actions.declineInvite();
    } else {
      actions.joinGroup(info.groupId, state.inviteToken ?? undefined);
    }
  };

  return (
    <div className="overlay" role="dialog" aria-modal="true" aria-label="Join group">
      <div className="card dialog">
        <h2>Join group &quot;{info.groupName}&quot;?</h2>
        <p className="muted">
          {info.memberCount} {info.memberCount === 1 ? "member" : "members"}. You will only be added if you press Join.
        </p>
        {info.alreadyMember && <p className="notice">You are already a member.</p>}
        <div className="row-buttons">
          <button className="secondary" onClick={actions.declineInvite}>
            Not now
          </button>
          <button onClick={onPrimary}>{info.alreadyMember ? "Open group" : "Join"}</button>
        </div>
      </div>
    </div>
  );
}
