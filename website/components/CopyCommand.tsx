"use client";

import { useEffect, useRef, useState } from "react";

const COPIED_FEEDBACK_MS = 1600;

export function CopyCommand({ command, prompt = "$" }: { command: string; prompt?: string }) {
  const [copied, setCopied] = useState(false);
  const resetTimer = useRef<ReturnType<typeof setTimeout>>(undefined);

  useEffect(() => () => clearTimeout(resetTimer.current), []);

  const copy = async () => {
    try {
      await navigator.clipboard.writeText(command);
    } catch {
      return;
    }
    setCopied(true);
    clearTimeout(resetTimer.current);
    resetTimer.current = setTimeout(() => setCopied(false), COPIED_FEEDBACK_MS);
  };

  return (
    <div className="command">
      <code>
        <span className="prompt" aria-hidden="true">{prompt} </span>
        {command}
      </code>
      <button type="button" className="ghost" onClick={copy} aria-live="polite">
        {copied ? "Copied" : "Copy"}
      </button>
    </div>
  );
}
