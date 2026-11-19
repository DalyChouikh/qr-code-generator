"use client";

import { useId, useState } from "react";
import type { FieldSpec } from "@/lib/fields";

export function TextField({
  spec,
  value,
  onChange,
}: {
  spec: FieldSpec;
  value: string;
  onChange: (value: string) => void;
}) {
  const id = useId();
  const [revealed, setRevealed] = useState(false);

  const common = {
    id,
    value,
    placeholder: spec.placeholder,
    required: spec.required,
    autoComplete: "off",
  };

  let control: React.ReactNode;
  if (spec.kind === "textarea") {
    control = <textarea {...common} rows={3} onChange={(e) => onChange(e.target.value)} />;
  } else if (spec.kind === "password") {
    control = (
      <div className="with-action">
        <input {...common} type={revealed ? "text" : "password"} onChange={(e) => onChange(e.target.value)} />
        <button type="button" className="ghost" onClick={() => setRevealed((r) => !r)}>
          {revealed ? "Hide" : "Show"}
        </button>
      </div>
    );
  } else {
    control = <input {...common} type={spec.kind} onChange={(e) => onChange(e.target.value)} />;
  }

  return (
    <div className="field">
      <label className="label" htmlFor={id}>
        {spec.label}
        {spec.required && <span className="required"> required</span>}
      </label>
      {control}
    </div>
  );
}
