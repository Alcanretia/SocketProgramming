import { useEffect, useState } from "react";
import QRCode from "qrcode";
import { useChat } from "../../socket/store";

/**
 * Shows a QR code for a group invite (special feature F2).
 * The QR only encodes a link: <this site>/?invite=<token>. Scanning it opens
 * the app, which asks the server what the invite is and shows "Join group X?".
 * Nobody is added until they press Join (R10).
 */
export function InviteQr({ groupId }: { groupId: string }) {
  const { state, actions } = useChat();
  const token = state.inviteCreated[groupId];
  const [dataUrl, setDataUrl] = useState<string | null>(null);
  const link = token ? `${location.origin}/?invite=${token}` : null;
  const localOnly = ["localhost", "127.0.0.1", "[::1]"].includes(location.hostname);

  useEffect(() => {
    if (!token) actions.createInvite(groupId);
  }, [groupId, token, actions]);

  useEffect(() => {
    if (!link) return;
    let cancelled = false;
    QRCode.toDataURL(link, { width: 220, margin: 1 })
      .then((url) => {
        if (!cancelled) setDataUrl(url);
      })
      .catch(() => {
        if (!cancelled) setDataUrl(null);
      });
    return () => {
      cancelled = true;
    };
  }, [link]);

  if (!link) return <p className="muted small">Creating invite...</p>;
  return (
    <div className="qr">
      {dataUrl && <img src={dataUrl} alt="QR code for the group invite link" width={220} height={220} />}
      <code className="small break">{link}</code>
      {localOnly && (
        <p className="notice small">
          This page was opened as <b>localhost</b>, which other devices cannot reach. Open the app with this computer&apos;s
          network address (for example http://192.168.x.x:8080) and create the QR again.
        </p>
      )}
    </div>
  );
}
