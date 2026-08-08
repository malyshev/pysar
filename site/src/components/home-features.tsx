import { FilePenLine, Languages, MicVocal, RefreshCw } from "lucide-react";
import { CornerTicks } from "@/components/corner-ticks";

const features = [
  {
    icon: FilePenLine,
    title: "From idea to a draft you can judge",
    body: "Start with a half-formed thought or a messy draft. Pysar helps you shape the piece step by step — and keeps earlier versions so you can compare.",
  },
  {
    icon: MicVocal,
    title: "Sounds like you, not a generic bot",
    body: "Capture how you write once. Later drafts follow that voice and structure — not a default AI tone.",
  },
  {
    icon: Languages,
    title: "Works in the editor you already use",
    body: "Claude Code, Cursor, or ChatGPT Codex — same writing help, in the place you already work.",
  },
  {
    icon: RefreshCw,
    title: "Real sources when a claim needs them",
    body: "When something should be checked, Pysar looks for real sources — not invented citations dressed up as research.",
  },
] as const;

/** Boxsi Features.tsx grid — dark framed section, Pysar claims. */
export function HomeFeatures() {
  return (
    <section
      id="features"
      className="size-full overflow-hidden bg-ink px-4 md:px-5 lg:px-7.5"
    >
      <div className="site-container relative z-10 border-white/10 py-16 lg:py-30">
        <div className="mb-12.5 px-5 md:px-12.5">
          <div className="grid gap-5 md:grid-cols-2">
            <div className="space-y-2.5 text-center md:text-start">
              <div className="inline-flex items-center justify-center rounded-full border border-white/15 px-4 py-1.5 text-xs font-semibold uppercase text-primary">
                What you get
              </div>
              <h2 className="text-[36px] text-primary-foreground md:text-[42px] lg:text-[52px]">
                Help that improves a piece{" "}
                <span className="text-white/40">without taking it over</span>
              </h2>
            </div>
            <div className="mx-auto max-w-62 text-center md:me-0 md:place-self-end md:text-end">
              Built around how{" "}
              <span className="text-primary-foreground">you</span> actually
              write — and when you say a piece is ready.
            </div>
          </div>
        </div>

        <div className="grid grid-cols-2 divide-x-2 divide-y-2 divide-white/10 border-y-2 border-white/10">
          <div className="relative p-5 md:p-7.5 lg:p-12.5">
            <div className="mb-12.5 flex size-10 items-center justify-center bg-white/10">
              <FilePenLine className="size-7 text-primary" aria-hidden="true" />
            </div>
            <h3 className="mb-2.5 text-xl text-primary-foreground">
              {features[0].title}
            </h3>
            <p className="text-white/50">{features[0].body}</p>
            <CornerTicks className="text-white/35" />
          </div>

          <div className="relative row-span-2 min-h-48 border-e-0 md:min-h-0 md:row-span-2">
            <div
              className="absolute inset-0 bg-[radial-gradient(circle_at_30%_20%,oklch(0.655_0.224_30.9_/_0.35),transparent_55%),repeating-linear-gradient(45deg,oklch(1_0_0_/_0.06)_0,oklch(1_0_0_/_0.06)_1px,transparent_0,transparent_50%)] bg-size-[auto,12px_12px]"
              aria-hidden="true"
            />
            <div className="relative flex h-full min-h-48 flex-col items-center justify-center gap-3 p-8 md:min-h-full">
              <span className="text-xs font-semibold uppercase tracking-wide text-primary">
                Your take stays yours
              </span>
              <p className="max-w-48 text-center text-sm text-primary-foreground/80">
                Voice, edits, and drafts live in your project — yours to shape.
                Nothing publishes until you do.
              </p>
            </div>
            <CornerTicks className="text-white/35" />
          </div>

          <div className="relative p-5 md:p-7.5 lg:p-12.5">
            <div className="mb-12.5 flex size-10 items-center justify-center bg-white/10">
              <MicVocal className="size-7 text-primary" aria-hidden="true" />
            </div>
            <h3 className="mb-2.5 text-xl text-primary-foreground">
              {features[1].title}
            </h3>
            <p className="text-white/50">{features[1].body}</p>
            <CornerTicks className="text-white/35" />
          </div>

          <div className="relative border-b-2 border-e-0 border-white/10 p-5 md:border-b-0 md:border-e-2 md:p-7.5 lg:p-12.5">
            <div className="mb-12.5 flex size-10 items-center justify-center bg-white/10">
              <Languages className="size-7 text-primary" aria-hidden="true" />
            </div>
            <h3 className="mb-2.5 text-xl text-primary-foreground">
              {features[2].title}
            </h3>
            <p className="text-white/50">{features[2].body}</p>
            <CornerTicks className="text-white/35" />
          </div>

          <div className="relative col-span-2 p-5 md:col-span-1 md:p-7.5 lg:p-12.5">
            <div className="mb-12.5 flex size-10 items-center justify-center bg-white/10">
              <RefreshCw className="size-7 text-primary" aria-hidden="true" />
            </div>
            <h3 className="mb-2.5 text-xl text-primary-foreground">
              {features[3].title}
            </h3>
            <p className="text-white/50">{features[3].body}</p>
            <CornerTicks className="text-white/35" />
          </div>
        </div>
      </div>
    </section>
  );
}
