"use client";

import { useMemo, useState } from "react";
import { FIELDS, INITIAL_FIELD_VALUES, WIFI_OPEN_NETWORK, type FieldValues } from "@/lib/fields";
import type { ContentTypeKey, OutputFormat } from "@/lib/qrgen";
import { useQrPreview } from "@/lib/useQrPreview";
import { ColorField } from "./ColorField";
import { SegmentedControl } from "./SegmentedControl";
import { TextField } from "./TextField";

const FORMAT_OPTIONS = [
  { value: "png", label: "PNG" },
  { value: "svg", label: "SVG" },
] as const;

const MIN_SIZE = 64;
const MAX_SIZE = 4096;

export function Generator({ intro }: { intro: React.ReactNode }) {
  const [type, setType] = useState<ContentTypeKey>("url");
  const [fieldValues, setFieldValues] = useState(INITIAL_FIELD_VALUES);
  const [format, setFormat] = useState<OutputFormat>("png");
  const [sizeInput, setSizeInput] = useState("256");
  const [fg, setFg] = useState("#000000");
  const [bg, setBg] = useState("#FFFFFF");

  const size = Number(sizeInput);
  const values: FieldValues = fieldValues[type];

  const request = useMemo(
    () => ({ type, fields: values, format, size, fg, bg }),
    [type, values, format, size, fg, bg],
  );
  const { api, preview, error } = useQrPreview(request);

  const setValue = (name: string, value: string) =>
    setFieldValues((prev) => ({ ...prev, [type]: { ...prev[type], [name]: value } }));

  const isOpenNetwork = type === "wifi" && values.encryption === WIFI_OPEN_NETWORK;
  const visibleFields = FIELDS[type].filter((f) => !(isOpenNetwork && f.name === "password"));

  return (
    <section className="generator" id="generate" aria-label="QR code generator">
      <div className="intro">{intro}</div>

      <form className="panel" onSubmit={(e) => e.preventDefault()}>
        <fieldset className="group">
          <legend>What goes in the code</legend>
          {api && (
            <SegmentedControl
              label="Content type"
              showLabel={false}
              variant="large"
              options={api.contentTypes.map((t) => ({ value: t.key, label: t.name, title: t.description }))}
              value={type}
              onChange={setType}
            />
          )}

          <div className="fields">
            {visibleFields.map((spec) => (
              <TextField
                key={`${type}-${spec.name}`}
                spec={spec}
                value={values[spec.name] ?? ""}
                onChange={(v) => setValue(spec.name, v)}
              />
            ))}

            {type === "wifi" && api && (
              <>
                <SegmentedControl
                  label="Security"
                  options={api.wifiEncryptionTypes.map((e) => ({ value: e.value, label: e.name }))}
                  value={values.encryption}
                  onChange={(v) => setValue("encryption", v)}
                />
                <label className="check">
                  <input
                    type="checkbox"
                    checked={values.hidden === "true"}
                    onChange={(e) => setValue("hidden", String(e.target.checked))}
                  />
                  Hidden network
                </label>
              </>
            )}
          </div>
        </fieldset>

        <fieldset className="group">
          <legend>How it looks</legend>
          <div className="look">
            <SegmentedControl label="Format" options={FORMAT_OPTIONS} value={format} onChange={setFormat} />
            <label className="field">
              <span className="label">Size in pixels</span>
              <input
                type="number"
                inputMode="numeric"
                min={MIN_SIZE}
                max={MAX_SIZE}
                step={1}
                value={sizeInput}
                onChange={(e) => setSizeInput(e.target.value)}
              />
            </label>
          </div>
          {api && (
            <>
              <ColorField label="Code color" value={fg} onChange={setFg} palette={api.colors} />
              <ColorField label="Background" value={bg} onChange={setBg} palette={api.colors} />
            </>
          )}
        </fieldset>
      </form>

      <div className="sheet-wrap">
        <figure className={preview ? "sheet ready" : "sheet"} style={{ background: bg }} aria-live="polite">
          {preview ? (
            <img src={preview.url} alt="Generated QR code" width={size} height={size} />
          ) : (
            <p className="sheet-message">{error ? toSentence(error) : "Starting the generator…"}</p>
          )}
        </figure>

        <div className="sheet-actions">
          <a
            className="button primary"
            href={preview?.url}
            download={`qrcode.${format}`}
            aria-disabled={!preview}
            onClick={(e) => !preview && e.preventDefault()}
          >
            Download {format.toUpperCase()}
          </a>
          <span className="meta">{Number.isFinite(size) ? `${size} × ${size} px` : "Enter a size"}</span>
        </div>

        <div className="sheet-notes">
          {preview && (
            <details className="payload">
              <summary>Show encoded text</summary>
              <pre>{preview.content}</pre>
            </details>
          )}
          <p className="privacy">Made in your browser by the same Go code the CLI runs. Nothing is uploaded.</p>
        </div>
      </div>
    </section>
  );
}

function toSentence(message: string) {
  const capitalized = message.charAt(0).toUpperCase() + message.slice(1);
  return capitalized.endsWith(".") ? capitalized : `${capitalized}.`;
}
