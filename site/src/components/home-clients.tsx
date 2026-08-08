import { CornerTicks } from "@/components/corner-ticks";

type HostLogo = {
  src: string;
  alt: string;
  /** Haft lists this host as experimental; Pysar has not shipped it yet. */
  planned?: boolean;
};

/**
 * Port of Boxsi Clients.tsx — same grid/paddings; official host logos.
 * Sources: site/public/hosts/SOURCES.txt
 */
const hosts: HostLogo[] = [
  { src: "/hosts/claude.svg", alt: "Claude Code" },
  { src: "/hosts/cursor.svg", alt: "Cursor" },
  { src: "/hosts/codex.svg", alt: "Codex" },
  { src: "/hosts/opencode.svg", alt: "OpenCode", planned: true },
];

export function HomeClients() {
  return (
    <section className="size-full overflow-hidden px-4 md:px-5 lg:px-7.5">
      <div className="relative z-10 mx-auto w-full max-w-250 border-x-2 border-zinc-100">
        <div className="grid grid-cols-2 divide-x-2 divide-y-2 divide-zinc-100 border-y-2 border-zinc-100 md:grid-cols-5">
          <div className="col-span-2 md:col-span-1">
            <div className="relative flex h-full flex-col items-center justify-center bg-zinc-100 px-7.5 py-9">
              <p className="text-sm text-zinc-900">
                Works with <span className="text-primary">Claude Code</span>,
                Cursor &amp; Codex
              </p>
              <CornerTicks className="text-zinc-500" />
            </div>
          </div>

          {hosts.map((host) => (
            <div
              key={host.alt}
              className="flex h-full flex-col items-center justify-center px-7.5 py-9"
            >
              {/* eslint-disable-next-line @next/next/no-img-element -- match template Image h-6.25 */}
              <img
                src={host.src}
                alt={host.planned ? `${host.alt} (planned)` : host.alt}
                title={host.planned ? `${host.alt} — planned` : undefined}
                loading="lazy"
                className={`h-6.25 w-auto ${host.planned ? "opacity-40" : ""}`}
              />
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
