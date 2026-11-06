import { useState, type FormEvent } from "react";
import { useChat } from "../../socket/store";

// ============================================================================
// P2 Frontend — งานย่อย 6.1: หน้ากรอกชื่อเข้าใช้งาน (NameEntry.tsx)
// ครอบคลุมเกณฑ์ R3 (ชื่อห้ามซ้ำ): ส่งคำสั่ง "register" ผ่าน WebSocket
// และแสดงข้อความแจ้งเตือน (state.registerError) ทันทีหากชื่อซ้ำหรือไม่ถูกต้อง
// ============================================================================

export function NameEntry() {
  const { state, actions } = useChat();
  const [name, setName] = useState(state.wantedName ?? "");

  const submit = (e: FormEvent) => {
    e.preventDefault();
    const trimmed = name.trim();
    if (!trimmed) return;
    actions.register(trimmed);
  };

  return (
    <form className="card name-entry" onSubmit={submit}>
      <h1>Socket Chat</h1>
      <p className="muted">Choose a unique name (1 to 20 characters) to join.</p>
      <div className="row">
        <input
          value={name}
          maxLength={20}
          autoFocus
          onChange={(e) => setName(e.target.value)}
          placeholder="Your name"
          aria-label="Your name"
        />
        <button type="submit" disabled={!name.trim() || state.status !== "open"}>
          Join
        </button>
      </div>
      {state.registerError && (
        <p className="error" role="alert">
          {state.registerError}
        </p>
      )}
    </form>
  );
}
