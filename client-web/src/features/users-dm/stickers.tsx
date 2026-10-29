// Sticker pack (special feature F3). Only the id travels over the WebSocket;
// the image is a static file served over HTTP from /stickers/<id>.svg.
// The ids must match server/stickers.go (a Go test checks the files exist).
// Images are original artwork made for this project.

export const STICKER_IDS = [
  "smile",
  "laugh",
  "love",
  "cry",
  "angry",
  "surprised",
  "sleepy",
  "thumbs_up",
] as const;

export function StickerView({ id, size = 96 }: { id: string; size?: number }) {
  // The id comes from the server's validated list, but encode it anyway.
  return (
    <img
      className="sticker"
      src={`/stickers/${encodeURIComponent(id)}.svg`}
      alt={`sticker: ${id}`}
      width={size}
      height={size}
    />
  );
}

export function StickerPicker({ onPick, disabled }: { onPick: (id: string) => void; disabled?: boolean }) {
  return (
    <div className="sticker-picker" role="listbox" aria-label="Stickers">
      {STICKER_IDS.map((id) => (
        <button key={id} type="button" className="sticker-btn" disabled={disabled} onClick={() => onPick(id)} title={id}>
          <StickerView id={id} size={44} />
        </button>
      ))}
    </div>
  );
}
