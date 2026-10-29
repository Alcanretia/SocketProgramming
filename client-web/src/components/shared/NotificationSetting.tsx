import { useState } from "react";
import { notifySupport, requestNotifyPermission, type NotifySupport } from "../../socket/notifications";

const TEXT: Record<NotifySupport, string> = {
  ready: "Desktop notifications are on.",
  "needs-permission": "Get a desktop notification when a message arrives while this tab is hidden.",
  denied: "Notifications are blocked in the browser settings. In-app alerts still work.",
  insecure: "Desktop notifications need HTTPS or localhost. In-app alerts, sound and the tab badge still work.",
  unsupported: "This browser has no Notification API. In-app alerts still work.",
};

export function NotificationSetting() {
  const [support, setSupport] = useState<NotifySupport>(notifySupport());
  return (
    <section>
      <h2>Notifications</h2>
      <p className="muted small">{TEXT[support]}</p>
      {support === "needs-permission" && (
        <button className="small" onClick={async () => setSupport(await requestNotifyPermission())}>
          Enable
        </button>
      )}
    </section>
  );
}
