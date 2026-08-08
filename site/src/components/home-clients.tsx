import { CornerTicks } from "@/components/corner-ticks";

type HostLogo = {
  src: string;
  alt: string;
  /** Visible product name when the asset is a mark-only glyph. */
  wordmark?: string;
  /** Not a shipped host yet — grayscale + Coming soon. */
  planned?: boolean;
};

/**
 * Port of Boxsi Clients.tsx — same grid/paddings; official host logos.
 * Sources: site/public/hosts/SOURCES.txt
 */
const hosts: HostLogo[] = [
  { src: "/hosts/claude.svg", alt: "Claude Code" },
  { src: "/hosts/cursor.svg", alt: "Cursor" },
  {
    src: "/hosts/chatgpt.svg",
    alt: "ChatGPT Codex",
    wordmark: "ChatGPT Codex",
  },
  { src: "/hosts/gemini.svg", alt: "Gemini", planned: true },
];

export function HomeClients() {
  return (
    <section className="size-full overflow-hidden px-4 md:px-5 lg:px-7.5">
      <div className="relative z-10 mx-auto w-full max-w-250 border-x-2 border-zinc-100">
        <div className="grid grid-cols-2 divide-x-2 divide-y-2 divide-zinc-100 border-y-2 border-zinc-100 md:grid-cols-5">
          <div className="col-span-2 md:col-span-1">
            <div className="relative flex h-full flex-col items-center justify-center bg-zinc-100 px-7.5 py-8">
              <p className="text-sm text-zinc-900">
                Works with <span className="text-primary">Claude Code</span>,
                Cursor &amp; ChatGPT Codex
              </p>
              <CornerTicks className="text-zinc-500" />
            </div>
          </div>

          {hosts.map((host) => (
            <div
              key={host.alt}
              className="flex h-full flex-col items-center justify-center gap-2 px-7.5 py-8"
            >
              {host.wordmark ? (
                <span
                  className={`inline-flex items-center gap-2.5 ${host.planned ? "grayscale opacity-40" : ""}`}
                >
                  {/* eslint-disable-next-line @next/next/no-img-element -- match template Image h-6.25 */}
                  <img
                    src={host.src}
                    alt=""
                    loading="lazy"
                    className="size-6.25 shrink-0"
                  />
                  <span className="text-sm font-semibold tracking-tight text-zinc-950">
                    {host.wordmark}
                  </span>
                </span>
              ) : (
                /* eslint-disable-next-line @next/next/no-img-element -- match template Image h-6.25 */
                <img
                  src={host.src}
                  alt={host.planned ? `${host.alt} (coming soon)` : host.alt}
                  loading="lazy"
                  className={`h-6.25 w-auto ${host.planned ? "grayscale opacity-40" : ""}`}
                />
              )}
              {host.planned ? (
                <span className="rounded-full border border-zinc-200 bg-zinc-50 px-2 py-0.5 text-[10px] font-semibold uppercase tracking-wide text-zinc-400">
                  Coming soon
                </span>
              ) : null}
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
