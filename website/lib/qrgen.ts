export type ContentTypeKey = "url" | "wifi" | "contact" | "email" | "sms" | "text";

export interface ContentType {
  key: ContentTypeKey;
  name: string;
  icon: string;
  description: string;
}

export interface NamedColor {
  name: string;
  hex: string;
}

export interface WifiEncryption {
  value: string;
  name: string;
}

export type OutputFormat = "png" | "svg";

export interface GenerateRequest {
  type: ContentTypeKey;
  fields: Record<string, string>;
  format: OutputFormat;
  size: number;
  fg: string;
  bg: string;
}

export type GenerateResult =
  | { data: Uint8Array; mime: string; content: string; error?: undefined }
  | { error: string };

export interface Qrgen {
  generate(request: GenerateRequest): GenerateResult;
  contentTypes: ContentType[];
  colors: NamedColor[];
  wifiEncryptionTypes: WifiEncryption[];
}

interface QrgenGlobal {
  generate(json: string): GenerateResult;
  contentTypes(): ContentType[];
  colors(): NamedColor[];
  wifiEncryptionTypes(): WifiEncryption[];
}

declare global {
  // Provided by Go's wasm_exec.js.
  class Go {
    importObject: WebAssembly.Imports;
    run(instance: WebAssembly.Instance): Promise<void>;
  }
  // eslint-disable-next-line no-var
  var qrgen: QrgenGlobal | undefined;
}

let loading: Promise<Qrgen> | null = null;

export function loadQrgen(): Promise<Qrgen> {
  loading ??= instantiate().catch((err) => {
    loading = null;
    throw err;
  });
  return loading;
}

async function instantiate(): Promise<Qrgen> {
  await loadScript("/wasm_exec.js");

  const go = new Go();
  const response = fetch("/qrgen.wasm");
  const { instance } = await WebAssembly.instantiateStreaming(response, go.importObject);
  void go.run(instance);

  const api = globalThis.qrgen;
  if (!api) throw new Error("The generator did not start.");

  return {
    generate: (request) => api.generate(JSON.stringify(request)),
    contentTypes: api.contentTypes(),
    colors: api.colors(),
    wifiEncryptionTypes: api.wifiEncryptionTypes(),
  };
}

function loadScript(src: string): Promise<void> {
  if (typeof Go !== "undefined") return Promise.resolve();
  return new Promise((resolve, reject) => {
    const script = document.createElement("script");
    script.src = src;
    script.onload = () => resolve();
    script.onerror = () => reject(new Error(`Could not load ${src}`));
    document.head.appendChild(script);
  });
}
