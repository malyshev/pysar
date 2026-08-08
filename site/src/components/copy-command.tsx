"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

type CopyCommandProps = {
  command: string;
  className?: string;
};

/** Soft command bar + filled Copy — layout from env-sentinel.dev (hero + docs). */
export function CopyCommand({ command, className }: CopyCommandProps) {
  const [copied, setCopied] = useState(false);
  const multiline = command.includes("\n");

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
    <div
      className={cn(
        "not-prose flex w-full gap-3 rounded-xl border border-zinc-200 bg-zinc-100/80 px-3 py-2.5 sm:px-4",
        multiline ? "items-start" : "items-center",
        className,
      )}
    >
      <code
        className={cn(
          "min-w-0 flex-1 text-start font-mono text-[13px] leading-relaxed text-zinc-800 sm:text-sm",
          multiline ? "whitespace-pre-wrap break-words" : "break-all",
        )}
      >
        {command}
      </code>
      <Button
        type="button"
        size="sm"
        onClick={onCopy}
        className={cn("shrink-0 rounded-lg px-4", multiline && "mt-0.5")}
      >
        {copied ? "Copied" : "Copy"}
      </Button>
    </div>
  );
}
