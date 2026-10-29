import { useState, type FormEvent } from "react";
import { useChat } from "../../socket/store";

export function NameEntry() {
  const { state, actions } = useChat();
  const [name, setName] = useState("");
  const connected = state.status === "open";

  const submit = (e: FormEvent) => {
    e.preventDefault();
    if (name.trim()) actions.register(name.trim());
  };

  return (
    <form className="card name-entry" onSubmit={submit}>
      <h1>Socket Chat</h1>
      <p className="muted">Pick a unique name to join.</p>
      <input
        autoFocus
        value={name}
        maxLength={20}
        onChange={(e) => setName(e.target.value)}
        placeholder="Your name"
        aria-label="Your name"
      />
      <button type="submit" disabled={!connected || !name.trim()}>
        {connected ? "Join" : "Connecting..."}
      </button>
      {state.registerError && <p className="error">{state.registerError}</p>}
    </form>
  );
}
