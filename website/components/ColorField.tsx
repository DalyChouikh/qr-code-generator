"use client";

import { useEffect, useId, useState } from "react";
import type { NamedColor } from "@/lib/qrgen";

const HEX = /^#?[0-9a-fA-F]{6}$/;

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
  const [draft, setDraft] = useState(value);
  useEffect(() => setDraft(value), [value]);

  const commit = (raw: string) => {
    setDraft(raw);
    if (HEX.test(raw.trim())) onChange(normalize(raw));
  };

  const invalid = !HEX.test(draft.trim());

  return (
    <div className="field color-field">
      <span className="label" id={`${id}-label`}>{label}</span>
      <div className="swatches" role="radiogroup" aria-labelledby={`${id}-label`}>
        {palette.map((c) => (
          <button
            key={c.name}
            type="button"
            role="radio"
            aria-checked={value.toUpperCase() === c.hex}
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
          onChange={(e) => onChange(normalize(e.target.value))}
        />
        <input
          type="text"
          aria-label={`${label}: hex value`}
          aria-invalid={invalid}
          spellCheck={false}
          maxLength={7}
          value={draft}
          onChange={(e) => commit(e.target.value)}
          onBlur={() => setDraft(value)}
        />
      </div>
    </div>
  );
}

function normalize(raw: string) {
  return `#${raw.trim().replace(/^#/, "").toUpperCase()}`;
}
