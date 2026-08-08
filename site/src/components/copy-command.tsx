"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";

/** Soft command bar + filled Copy — layout from env-sentinel.dev. */
export function CopyCommand({ command }: { command: string }) {
  const [copied, setCopied] = useState(false);

  async function onCopy() {
    try {
      await navigator.clipboard.writeText(command);
      setCopied(true);
      window.setTimeout(() => setCopied(false), 2000);
    } catch {
      // Clipboard can fail in non-secure contexts — leave UI unchanged.
    }
  }

  return (
    <div className="flex w-full items-center gap-3 rounded-xl border border-zinc-200 bg-zinc-100/80 px-3 py-2.5 sm:px-4">
      <code className="min-w-0 flex-1 break-all text-start font-mono text-[13px] leading-relaxed text-zinc-800 sm:text-sm">
        {command}
      </code>
      <Button
        type="button"
        size="sm"
        onClick={onCopy}
        className="shrink-0 rounded-lg px-4"
      >
        {copied ? "Copied" : "Copy"}
      </Button>
    </div>
  );
}
