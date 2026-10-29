import type { ServerFrame } from "./types";

export type ConnStatus = "connecting" | "open" | "closed";

interface Handlers {
  onFrame: (f: ServerFrame) => void;
  onStatus: (s: ConnStatus) => void;
}

/**
 * Thin wrapper over the browser's plain WebSocket API.
 *  - one JSON envelope per text frame: {type, payload}
 *  - reconnects with backoff after a drop (the server forgets us when the
 *    socket closes, so the store re-registers the name once we are open again)
 */
export class Connection {
  private ws: WebSocket | null = null;
  private retries = 0;
  private timer: number | undefined;
  private stopped = false;

  constructor(
    private readonly url: string,
    private readonly handlers: Handlers,
  ) {}

  static defaultUrl(): string {
    const scheme = location.protocol === "https:" ? "wss" : "ws";
    return `${scheme}://${location.host}/ws`;
  }

  connect(): void {
    this.stopped = false;
    this.handlers.onStatus("connecting");
    const ws = new WebSocket(this.url);
    this.ws = ws;

    ws.onopen = () => {
      this.retries = 0;
      this.handlers.onStatus("open");
    };
    ws.onmessage = (e) => {
      try {
        const frame = JSON.parse(String(e.data)) as ServerFrame;
        if (frame && typeof frame.type === "string") this.handlers.onFrame(frame);
      } catch {
        // ignore a frame we cannot parse; the server never sends one
      }
    };
    ws.onclose = () => {
      if (this.ws !== ws) return; // an older socket closing after a reconnect
      this.handlers.onStatus("closed");
      if (this.stopped) return;
      const delay = Math.min(10_000, 500 * 2 ** this.retries++);
      this.timer = window.setTimeout(() => this.connect(), delay);
    };
    ws.onerror = () => ws.close();
  }

  send(type: string, payload: Record<string, unknown> = {}): boolean {
    if (this.ws?.readyState !== WebSocket.OPEN) return false;
    this.ws.send(JSON.stringify({ type, payload }));
    return true;
  }

  close(): void {
    this.stopped = true;
    window.clearTimeout(this.timer);
    const ws = this.ws;
    this.ws = null;
    ws?.close();
  }
}
