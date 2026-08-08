import Link from "next/link";
import { CornerTicks } from "@/components/corner-ticks";
import { siteConfig } from "@/lib/site";

/** Boxsi Clients.tsx strip — honest cells, no fake logos. */
export function HomeClients() {
  const cells = [
    { label: "Docs", href: "/docs" },
    { label: "Install", href: "/docs/install" },
    { label: "GitHub", href: siteConfig.github, external: true },
    { label: "MCP", href: "/docs/mcp-and-skills" },
  ] as const;

  return (
    <section className="size-full overflow-hidden px-4 md:px-5 lg:px-7.5">
      <div className="site-container relative z-10">
        <div className="grid grid-cols-2 divide-x-2 divide-y-2 divide-frame border-y-2 border-frame md:grid-cols-5">
          <div className="col-span-2 md:col-span-1">
            <div className="relative flex h-full flex-col items-center justify-center bg-muted px-7.5 py-9">
              <p className="text-sm text-ink">
                Open source · <span className="text-primary">MIT</span> · CLI +
                MCP
              </p>
              <CornerTicks />
            </div>
          </div>

          {cells.map((cell) => (
            <div
              key={cell.label}
              className="flex h-full flex-col items-center justify-center px-7.5 py-9"
            >
              {"external" in cell && cell.external ? (
                <a
                  href={cell.href}
                  className="text-sm font-semibold uppercase tracking-wide text-ink transition-colors hover:text-primary"
                  rel="noopener noreferrer"
                >
                  {cell.label}
                </a>
              ) : (
                <Link
                  href={cell.href}
                  className="text-sm font-semibold uppercase tracking-wide text-ink transition-colors hover:text-primary"
                >
                  {cell.label}
                </Link>
              )}
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
