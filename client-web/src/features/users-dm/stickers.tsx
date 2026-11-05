// ============================================================================
// P2 Frontend — งานย่อย 6.3: ระบบเลือกและแสดงผลสติกเกอร์ (stickers.tsx)
// ครอบคลุมเกณฑ์ Special Feature F3 (Default Stickers):
//   - รายการรหัสสติกเกอร์ทั้ง 8 ตัวต้องตรงกับ stickerPack ใน server/stickers.go
//   - ส่งเฉพาะรหัส id ผ่าน WebSocket ส่วนรูปภาพดึงจากไฟล์ Static (/stickers/<id>.svg)
// ============================================================================

export interface StickerDef {
  id: string;
  label: string;
}

export const STICKERS: readonly StickerDef[] = [
  { id: "smile", label: "Smile" },
  { id: "laugh", label: "Laugh" },
  { id: "love", label: "Love" },
  { id: "cry", label: "Cry" },
  { id: "angry", label: "Angry" },
  { id: "surprised", label: "Surprised" },
  { id: "sleepy", label: "Sleepy" },
  { id: "thumbs_up", label: "Thumbs up" },
];

/** แสดงรูปสติกเกอร์ในหน้าต่างแชท (พร้อมตรวจสอบก่อนว่ามีรหัสอยู่ในระบบจริง) */
export function StickerView({ id }: { id: string }) {
  const known = STICKERS.find((s) => s.id === id);
  if (!known) {
    return <span className="bubble">[sticker: {id}]</span>;
  }
  return (
    <img
      className="sticker"
      src={`/stickers/${known.id}.svg`}
      alt={known.label}
      title={known.label}
      width={96}
      height={96}
    />
  );
}

/** แผงปุ่มกดเลือกสติกเกอร์ทั้ง 8 แบบข้างกล่องพิมพ์ข้อความ */
export function StickerPicker({ onPick }: { onPick: (id: string) => void }) {
  return (
    <div className="sticker-picker" role="listbox" aria-label="Pick a sticker">
      {STICKERS.map((s) => (
        <button
          key={s.id}
          type="button"
          className="sticker-btn"
          onClick={() => onPick(s.id)}
          title={s.label}
          aria-label={`Send ${s.label} sticker`}
        >
          <img src={`/stickers/${s.id}.svg`} alt={s.label} width={56} height={56} />
        </button>
      ))}
    </div>
  );
}
