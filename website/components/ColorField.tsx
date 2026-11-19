"use client";

import { useId, useState } from "react";
import type { NamedColor } from "@/lib/qrgen";

const HEX_COLOR = /^#?[0-9a-fA-F]{6}$/;

export function ColorField({
  label,
  value,
  onChange,
  palette,
}: {
  label: string;
  value: string;
  onChange: (hex: string) => void;
  palette: NamedColor[];
}) {
  const id = useId();
  // Holds partial hex input while typing; null means "show the committed value".
  const [draft, setDraft] = useState<string | null>(null);
  const shown = draft ?? value;

  const editHex = (raw: string) => {
    setDraft(raw);
    if (HEX_COLOR.test(raw.trim())) onChange(normalizeHex(raw));
  };

  return (
    <div className="field color-field">
      <span className="label" id={`${id}-label`}>{label}</span>
      <div className="swatches" role="radiogroup" aria-labelledby={`${id}-label`}>
        {palette.map((c) => (
          <button
            key={c.name}
            type="button"
            role="radio"
            aria-checked={value === c.hex}
            aria-label={c.name}
            title={c.name}
            className="swatch"
            style={{ background: c.hex }}
            onClick={() => onChange(c.hex)}
          />
        ))}
      </div>
      <div className="hex">
        <input
          type="color"
          aria-label={`${label}: pick any color`}
          value={value.toLowerCase()}
          onChange={(e) => onChange(normalizeHex(e.target.value))}
        />
        <input
          type="text"
          aria-label={`${label}: hex value`}
          aria-invalid={!HEX_COLOR.test(shown.trim())}
          spellCheck={false}
          maxLength={7}
          value={shown}
          onChange={(e) => editHex(e.target.value)}
          onBlur={() => setDraft(null)}
        />
      </div>
    </div>
  );
}

function normalizeHex(raw: string) {
  return `#${raw.trim().replace(/^#/, "").toUpperCase()}`;
}
