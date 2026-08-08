import { ChevronRight } from "lucide-react";
import Link from "next/link";
import { CornerTicks } from "@/components/corner-ticks";

type Node = {
  src: string;
  alt: string;
  /** Percent positions matching Boxsi integration.svg node plates. */
  left: string;
  top: string;
  planned?: boolean;
};

/** Node plates from Boxsi integration.svg (896×176 viewBox → %). */
const nodes: Node[] = [
  { src: "/hosts/claude.svg", alt: "Claude Code", left: "0%", top: "8.5%" },
  { src: "/hosts/cursor.svg", alt: "Cursor", left: "10.7%", top: "8.5%" },
  {
    src: "/hosts/chatgpt.svg",
    alt: "ChatGPT Codex",
    left: "21.4%",
    top: "8.5%",
  },
  {
    src: "/hosts/gemini-icon.svg",
    alt: "Gemini",
    left: "0%",
    top: "63%",
    planned: true,
  },
  { src: "/hosts/claude.svg", alt: "Claude Code", left: "10.7%", top: "63%" },
  {
    src: "/hosts/chatgpt.svg",
    alt: "ChatGPT Codex",
    left: "73%",
    top: "8.5%",
  },
  { src: "/hosts/cursor.svg", alt: "Cursor", left: "83.7%", top: "8.5%" },
  { src: "/hosts/claude.svg", alt: "Claude Code", left: "94.4%", top: "8.5%" },
  {
    src: "/hosts/gemini-icon.svg",
    alt: "Gemini",
    left: "94.4%",
    top: "63%",
    planned: true,
  },
  {
    src: "/hosts/chatgpt.svg",
    alt: "ChatGPT Codex",
    left: "83.7%",
    top: "63%",
  },
];

/**
 * Port of Boxsi Integrations.tsx — same header/CTA rhythm + connector frame.
 * Third-party SaaS icons replaced with official host logos.
 */
export function HomeIntegrations() {
  return (
    <section className="size-full overflow-hidden px-4 md:px-5 lg:px-7.5">
      <div className="relative z-10 mx-auto w-full max-w-250 border-x-2 border-zinc-100 py-16 lg:py-30">
        <div className="mx-auto mb-12.5 max-w-135 px-12.5">
          <div className="mb-10 space-y-2.5 text-center">
            <div className="inline-flex items-center justify-center rounded-full border-2 border-zinc-200 px-4 py-1.5 text-xs font-semibold uppercase text-primary">
              Seamlessly Connected
            </div>

            <h2 className="mb-2.5 text-section md:text-section-md lg:text-section-lg">
              Works in editors{" "}
              <span className="text-zinc-500">you already use</span>
            </h2>

            <div className="relative inline-flex p-1.5">
              <Link
                href="/docs/install"
                className="inline-flex items-center justify-center gap-2 rounded bg-primary px-4 py-3 text-xs font-medium uppercase text-white transition-all duration-500 hover:bg-primary-hover"
              >
                Install
                <ChevronRight className="size-4" aria-hidden="true" />
              </Link>
              <CornerTicks className="text-zinc-500" />
            </div>
          </div>
        </div>

        <div className="px-7.5">
          <div className="relative mx-auto w-full max-w-full">
            {/* eslint-disable-next-line @next/next/no-img-element */}
            <img
              src="/patterns/integration-frame.svg"
              alt=""
              className="h-auto w-full"
              loading="lazy"
            />

            {nodes.map((node, i) => (
              <div
                key={`${node.alt}-${i}`}
                className="absolute flex size-10 items-center justify-center rounded-sm bg-zinc-100 p-1.5 md:size-12.5"
                style={{ left: node.left, top: node.top }}
                title={node.planned ? `${node.alt} — Coming soon` : undefined}
              >
                {/* eslint-disable-next-line @next/next/no-img-element */}
                <img
                  src={node.src}
                  alt={node.planned ? `${node.alt} (coming soon)` : node.alt}
                  className={`h-auto max-h-full w-full object-contain ${node.planned ? "grayscale opacity-40" : ""}`}
                  loading="lazy"
                />
              </div>
            ))}

            <div
              className="absolute top-1/2 left-1/2 flex size-16 -translate-x-1/2 -translate-y-1/2 items-center justify-center rounded border-2 border-zinc-200 bg-zinc-100 md:size-20"
              aria-hidden="true"
            >
              {/* eslint-disable-next-line @next/next/no-img-element */}
              <img
                src="/logo-mark.svg"
                alt=""
                className="pixelated size-logo-mark"
              />
              <CornerTicks className="text-zinc-500" />
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
