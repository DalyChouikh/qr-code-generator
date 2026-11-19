"use client";

import { useEffect, useState } from "react";
import { loadQrgen, type GenerateRequest, type Qrgen } from "./qrgen";

const DEBOUNCE_MS = 120;

export interface QrPreview {
  url: string;
  content: string;
}

export interface QrPreviewState {
  api: Qrgen | null;
  preview: QrPreview | null;
  error: string | null;
}

/** Loads the WASM generator and re-renders the code (debounced) whenever the request changes. */
export function useQrPreview(request: GenerateRequest): QrPreviewState {
  const [api, setApi] = useState<Qrgen | null>(null);
  const [preview, setPreview] = useState<QrPreview | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadQrgen().then(setApi, (err: Error) =>
      setError(`The generator failed to load (${err.message}). Reload the page to try again`),
    );
  }, []);

  useEffect(() => {
    if (!api) return;
    const timer = setTimeout(() => {
      const result = api.generate(request);
      if (result.error !== undefined) {
        setError(result.error);
        setPreview(null);
        return;
      }
      setError(null);
      const blob = new Blob([result.data as BlobPart], { type: result.mime });
      setPreview({ url: URL.createObjectURL(blob), content: result.content });
    }, DEBOUNCE_MS);
    return () => clearTimeout(timer);
  }, [api, request]);

  useEffect(() => {
    if (!preview) return;
    return () => URL.revokeObjectURL(preview.url);
  }, [preview]);

  return { api, preview, error };
}
