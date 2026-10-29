import { useEffect, useRef, useState, type FormEvent, type ReactNode } from "react";
import type { ChatMessage } from "../../socket/types";
import { StickerPicker, StickerView } from "../../features/users-dm/stickers";

interface Props {
  title: string;
  subtitle?: string;
  me: string;
  messages: ChatMessage[];
  disabledReason?: string; // when set, the composer is locked and this is shown
  onSendText: (text: string) => void;
  onSendSticker: (id: string) => void;
  header?: ReactNode; // e.g. the AI summary card
  onBack?: () => void;
}

/** One chat room: a message window plus a chat box (R6). Used for DMs and groups. */
export function ChatRoom({ title, subtitle, me, messages, disabledReason, onSendText, onSendSticker, header, onBack }: Props) {
  const [text, setText] = useState("");
  const [showStickers, setShowStickers] = useState(false);
  const end = useRef<HTMLDivElement>(null);
  const locked = Boolean(disabledReason);

  useEffect(() => {
    end.current?.scrollIntoView({ block: "end" });
  }, [messages.length, title]);

  const submit = (e: FormEvent) => {
    e.preventDefault();
    const t = text.trim();
    if (!t || locked) return;
    onSendText(t);
    setText("");
  };

  return (
    <div className="chatroom">
      <header className="chat-header">
        {onBack && (
          <button className="back" onClick={onBack} aria-label="Back to lists">
            ‹
          </button>
        )}
        <div>
          <strong>{title}</strong>
          {subtitle && <div className="muted small">{subtitle}</div>}
        </div>
      </header>

      {header}

      <div className="messages" role="log" aria-live="polite">
        {messages.length === 0 && <p className="muted center">No messages yet.</p>}
        {messages.map((m) => {
          const mine = m.from.toLowerCase() === me.toLowerCase();
          return (
            <div key={m.msgId} className={`msg ${mine ? "mine" : ""}`}>
              {!mine && <div className="from">{m.from}</div>}
              {m.kind === "sticker" && m.stickerId ? (
                <StickerView id={m.stickerId} />
              ) : (
                // React escapes this text; never use dangerouslySetInnerHTML with user content.
                <div className="bubble">{m.text}</div>
              )}
              <time className="muted small">{new Date(m.ts).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })}</time>
            </div>
          );
        })}
        <div ref={end} />
      </div>

      {showStickers && !locked && (
        <StickerPicker
          onPick={(id) => {
            onSendSticker(id);
            setShowStickers(false);
          }}
        />
      )}

      {disabledReason && <p className="notice">{disabledReason}</p>}
      <form className="composer" onSubmit={submit}>
        <button type="button" className="icon" disabled={locked} onClick={() => setShowStickers((v) => !v)} aria-label="Stickers" title="Stickers">
          ☺
        </button>
        <input
          value={text}
          maxLength={1000}
          disabled={locked}
          onChange={(e) => setText(e.target.value)}
          placeholder="Type a message"
          aria-label="Message"
        />
        <button type="submit" disabled={locked || !text.trim()}>
          Send
        </button>
      </form>
    </div>
  );
}
