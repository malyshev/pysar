import { cn } from "@/lib/utils";

type SiteLogoProps = {
  variant?: "light" | "dark";
  className?: string;
  /** Mark size override; defaults to size-logo-mark. */
  markClassName?: string;
  priority?: boolean;
};

/**
 * Template lockup: mark + bold sans wordmark.
 * Tokens: text-logo, tracking-logo, size-logo-mark (see globals @theme).
 */
export function SiteLogo({
  variant = "light",
  className,
  markClassName,
  priority = false,
}: SiteLogoProps) {
  const markSrc =
    variant === "dark" ? "/logo-mark-dark.svg" : "/logo-mark.svg";

  return (
    <span className={cn("inline-flex items-center gap-2.5", className)}>
      {/* eslint-disable-next-line @next/next/no-img-element -- next/image softens pixel stitches */}
      <img
        src={markSrc}
        alt=""
        width={304}
        height={304}
        className={cn("pixelated size-logo-mark shrink-0", markClassName)}
        decoding={priority ? "sync" : "async"}
        fetchPriority={priority ? "high" : "auto"}
      />
      <span
        className={cn(
          "text-logo font-bold uppercase leading-none tracking-logo",
          variant === "dark" ? "text-white" : "text-zinc-950",
        )}
      >
        Pysar
      </span>
    </span>
  );
}
