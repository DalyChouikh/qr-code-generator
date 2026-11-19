import type { ContentTypeKey } from "./qrgen";
import { repoUrl } from "./site";

export type FieldKind = "text" | "textarea" | "email" | "tel" | "url" | "password";

export interface FieldSpec {
  name: string;
  label: string;
  kind: FieldKind;
  placeholder?: string;
  required?: boolean;
}

export type FieldValues = Record<string, string>;

// Field names are the keys cmd/qrgen-wasm reads; the sets follow internal/ui/template_wizard.go.
export const FIELDS: Record<ContentTypeKey, FieldSpec[]> = {
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

export const WIFI_OPEN_NETWORK = "nopass";

export const INITIAL_FIELD_VALUES: Record<ContentTypeKey, FieldValues> = {
  url: { content: repoUrl },
  text: {},
  wifi: { encryption: "WPA" },
  contact: {},
  email: {},
  sms: {},
};
