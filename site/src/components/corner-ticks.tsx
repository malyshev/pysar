import { cn } from "@/lib/utils";

/** Geometric corner ticks used on framed logo/CTA blocks. */
export function CornerTicks({ className }: { className?: string }) {
  return (
    <div className={cn("decoration", className)} aria-hidden="true">
      <span className="decoration-top" />
      <span className="decoration-bottom" />
    </div>
  );
}
