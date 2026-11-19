"use client";

import { useState } from "react";

export function CopyCommand({ command, prompt = "$" }: { command: string; prompt?: string }) {
  const [copied, setCopied] = useState(false);

  const copy = async () => {
    try {
      await navigator.clipboard.writeText(command);
      setCopied(true);
      setTimeout(() => setCopied(false), 1600);
    } catch {
      setCopied(false);
    }
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
