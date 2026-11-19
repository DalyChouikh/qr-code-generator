"use client";

import { useEffect, useMemo, useState } from "react";
import {
  loadQrgen,
  type ContentTypeKey,
  type OutputFormat,
  type Qrgen,
} from "@/lib/qrgen";
import { repoUrl } from "@/lib/site";
import { ColorField } from "./ColorField";

type FieldKind = "text" | "textarea" | "email" | "tel" | "url" | "password";

interface FieldSpec {
  name: string;
  label: string;
  kind: FieldKind;
  placeholder?: string;
  required?: boolean;
}

// Field sets follow internal/ui/template_wizard.go so the web form asks for exactly what the TUI asks for.
const FIELDS: Record<ContentTypeKey, FieldSpec[]> = {
  url: [{ name: "content", label: "Link", kind: "url", placeholder: "https://example.com", required: true }],
  text: [{ name: "content", label: "Text", kind: "textarea", placeholder: "Anything you want to share", required: true }],
  wifi: [
    { name: "ssid", label: "Network name (SSID)", kind: "text", placeholder: "MyNetwork", required: true },
    { name: "password", label: "Password", kind: "password" },
  ],
  contact: [
    { name: "firstName", label: "First name", kind: "text", placeholder: "Ada" },
    { name: "lastName", label: "Last name", kind: "text", placeholder: "Lovelace" },
    { name: "phone", label: "Phone", kind: "tel", placeholder: "+1 555 0100" },
    { name: "email", label: "Email", kind: "email", placeholder: "ada@example.com" },
    { name: "organization", label: "Organization", kind: "text" },
    { name: "title", label: "Job title", kind: "text" },
    { name: "url", label: "Website", kind: "url" },
  ],
  email: [
    { name: "address", label: "To", kind: "email", placeholder: "hello@example.com", required: true },
    { name: "subject", label: "Subject", kind: "text" },
    { name: "body", label: "Message", kind: "textarea" },
  ],
  sms: [
    { name: "phone", label: "Phone number", kind: "tel", placeholder: "+1 555 0100", required: true },
    { name: "message", label: "Message", kind: "textarea" },
  ],
};

const MIN_SIZE = 64;
const MAX_SIZE = 4096;

interface Output {
  url: string;
  content: string;
}

export function Generator({ intro }: { intro: React.ReactNode }) {
  const [api, setApi] = useState<Qrgen | null>(null);
  const [loadError, setLoadError] = useState<string | null>(null);

  const [type, setType] = useState<ContentTypeKey>("url");
  const [fields, setFields] = useState<Record<ContentTypeKey, Record<string, string>>>({
    url: { content: repoUrl },
    text: {},
    wifi: { encryption: "WPA" },
    contact: {},
    email: {},
    sms: {},
  });
  const [format, setFormat] = useState<OutputFormat>("png");
  const [sizeInput, setSizeInput] = useState("256");
  const [fg, setFg] = useState("#000000");
  const [bg, setBg] = useState("#FFFFFF");
  const [showPassword, setShowPassword] = useState(false);

  const [output, setOutput] = useState<Output | null>(null);
  const [problem, setProblem] = useState<string | null>(null);

  useEffect(() => {
    loadQrgen().then(setApi, (err: Error) => setLoadError(err.message));
  }, []);

  const size = Number(sizeInput);
  const current = fields[type];

  const request = useMemo(
    () => ({ type, fields: current, format, size, fg, bg }),
    [type, current, format, size, fg, bg],
  );

  useEffect(() => {
    if (!api) return;
    const timer = setTimeout(() => {
      const result = api.generate(request);
      if (result.error !== undefined) {
        setProblem(result.error);
        setOutput(null);
        return;
      }
      setProblem(null);
      const blob = new Blob([result.data as BlobPart], { type: result.mime });
      setOutput({ url: URL.createObjectURL(blob), content: result.content });
    }, 120);
    return () => clearTimeout(timer);
  }, [api, request]);

  useEffect(() => {
    if (!output) return;
    return () => URL.revokeObjectURL(output.url);
  }, [output]);

  const setField = (name: string, value: string) =>
    setFields((prev) => ({ ...prev, [type]: { ...prev[type], [name]: value } }));

  const types = api?.contentTypes ?? [];
  const encryption = current.encryption ?? "WPA";

  return (
    <section className="generator" id="generate" aria-label="QR code generator">
      <div className="intro">{intro}</div>
      <form className="panel" onSubmit={(e) => e.preventDefault()}>
        <fieldset className="group">
          <legend>What goes in the code</legend>
          <div className="types" role="radiogroup" aria-label="Content type">
            {types.map((t) => (
              <button
                key={t.key}
                type="button"
                role="radio"
                aria-checked={type === t.key}
                className="type"
                onClick={() => setType(t.key)}
                title={t.description}
              >
                {t.name}
              </button>
            ))}
          </div>

          <div className="fields">
            {FIELDS[type].map((f) => (
              <Field
                key={`${type}-${f.name}`}
                spec={f}
                value={current[f.name] ?? ""}
                onChange={(v) => setField(f.name, v)}
                hidden={f.name === "password" && encryption === "nopass"}
                revealed={showPassword}
                onToggleReveal={() => setShowPassword((s) => !s)}
              />
            ))}

            {type === "wifi" && api && (
              <>
                <div className="field">
                  <span className="label" id="enc-label">Security</span>
                  <div className="segmented" role="radiogroup" aria-labelledby="enc-label">
                    {api.wifiEncryptionTypes.map((e) => (
                      <button
                        key={e.value}
                        type="button"
                        role="radio"
                        aria-checked={encryption === e.value}
                        onClick={() => setField("encryption", e.value)}
                      >
                        {e.name}
                      </button>
                    ))}
                  </div>
                </div>
                <label className="check">
                  <input
                    type="checkbox"
                    checked={current.hidden === "true"}
                    onChange={(e) => setField("hidden", String(e.target.checked))}
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
            <div className="field">
              <span className="label" id="format-label">Format</span>
              <div className="segmented" role="radiogroup" aria-labelledby="format-label">
                {(["png", "svg"] as const).map((f) => (
                  <button key={f} type="button" role="radio" aria-checked={format === f} onClick={() => setFormat(f)}>
                    {f.toUpperCase()}
                  </button>
                ))}
              </div>
            </div>
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
        <figure className={output ? "sheet ready" : "sheet"} style={{ background: bg }} aria-live="polite">
          {output ? (
            // Blob URLs can't go through next/image.
            // eslint-disable-next-line @next/next/no-img-element
            <img src={output.url} alt="Generated QR code" width={size} height={size} />
          ) : (
            <p className="sheet-message">
              {loadError
                ? `The generator failed to load: ${loadError}. Reload the page to try again.`
                : problem
                  ? capitalize(problem)
                  : "Starting the generator…"}
            </p>
          )}
        </figure>

        <div className="sheet-actions">
          <a
            className="button primary"
            href={output?.url}
            download={`qrcode.${format}`}
            aria-disabled={!output}
            onClick={(e) => !output && e.preventDefault()}
          >
            Download {format.toUpperCase()}
          </a>
          <span className="meta">
            {Number.isFinite(size) ? `${size} × ${size} px` : "Enter a size"}
          </span>
        </div>

        <div className="sheet-notes">
          {output && (
            <details className="payload">
              <summary>Show encoded text</summary>
              <pre>{output.content}</pre>
            </details>
          )}
          <p className="privacy">Made in your browser by the same Go code the CLI runs. Nothing is uploaded.</p>
        </div>
      </div>
    </section>
  );
}

function Field({
  spec,
  value,
  onChange,
  hidden,
  revealed,
  onToggleReveal,
}: {
  spec: FieldSpec;
  value: string;
  onChange: (v: string) => void;
  hidden: boolean;
  revealed: boolean;
  onToggleReveal: () => void;
}) {
  if (hidden) return null;
  const label = (
    <span className="label">
      {spec.label}
      {spec.required && <span className="required"> required</span>}
    </span>
  );

  if (spec.kind === "textarea") {
    return (
      <label className="field">
        {label}
        <textarea rows={3} value={value} placeholder={spec.placeholder} onChange={(e) => onChange(e.target.value)} />
      </label>
    );
  }

  if (spec.kind === "password") {
    return (
      <div className="field">
        <label htmlFor="wifi-password">{label}</label>
        <div className="with-action">
          <input
            id="wifi-password"
            type={revealed ? "text" : "password"}
            autoComplete="off"
            value={value}
            onChange={(e) => onChange(e.target.value)}
          />
          <button type="button" className="ghost" onClick={onToggleReveal}>
            {revealed ? "Hide" : "Show"}
          </button>
        </div>
      </div>
    );
  }

  return (
    <label className="field">
      {label}
      <input
        type={spec.kind}
        value={value}
        placeholder={spec.placeholder}
        autoComplete="off"
        onChange={(e) => onChange(e.target.value)}
      />
    </label>
  );
}

function capitalize(s: string) {
  return s.charAt(0).toUpperCase() + s.slice(1) + (s.endsWith(".") ? "" : ".");
}
