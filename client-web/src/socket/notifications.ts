// Notifications (special feature F4).
//
// Which layer does what:
//   - the SERVER only delivers the message (and keeps the unread counts);
//   - THIS file runs in the browser and decides how to alert the user.
// The browser Notification API needs a secure context (HTTPS or localhost).
// On a plain-HTTP LAN demo it is unavailable, so every alert also has an
// in-app fallback (toast, title badge, sound) that always works.

export type NotifySupport = "ready" | "needs-permission" | "denied" | "insecure" | "unsupported";

export function notifySupport(): NotifySupport {
  if (typeof Notification === "undefined") return "unsupported";
  if (!window.isSecureContext) return "insecure";
  if (Notification.permission === "granted") return "ready";
  if (Notification.permission === "denied") return "denied";
  return "needs-permission";
}

/** Must be called from a click handler (browsers block unprompted requests). */
export async function requestNotifyPermission(): Promise<NotifySupport> {
  if (notifySupport() === "needs-permission") {
    try {
      await Notification.requestPermission();
    } catch {
      // some browsers throw instead of resolving
    }
  }
  return notifySupport();
}

export function showSystemNotification(title: string, body: string, onClick: () => void): void {
  if (notifySupport() !== "ready") return;
  try {
    const n = new Notification(title, { body, tag: title });
    n.onclick = () => {
      window.focus();
      onClick();
      n.close();
    };
  } catch {
    // e.g. mobile Chrome requires a service worker; the in-app fallback still runs
  }
}

let audio: AudioContext | null = null;

export function beep(): void {
  try {
    audio ??= new AudioContext();
    const osc = audio.createOscillator();
    const gain = audio.createGain();
    osc.frequency.value = 880;
    gain.gain.value = 0.04;
    osc.connect(gain).connect(audio.destination);
    osc.start();
    osc.stop(audio.currentTime + 0.12);
  } catch {
    // audio can be blocked until the user interacts with the page
  }
}

const BASE_TITLE = "Socket Chat";

export function setTitleBadge(unread: number): void {
  document.title = unread > 0 ? `(${unread}) ${BASE_TITLE}` : BASE_TITLE;
}
