import { Circle } from "lucide-react";
import { SiteShell } from "@/components/site-shell";

/**
 * Boxsi About-style docs frame: max-w-250 rail + slots for nav + reading column.
 * Matches home/footer geometry (dec / template language), not site-container framed shell.
 */
export function DocsLayout({
  nav,
  children,
}: {
  nav: React.ReactNode;
  children: React.ReactNode;
}) {
  return (
    <SiteShell framed={false}>
      <section className="size-full overflow-hidden px-4 md:px-5 lg:px-7.5">
        <div className="relative z-10 mx-auto w-full max-w-250 border-x-2 border-zinc-100 pt-32 pb-16 md:pt-37.5 lg:pb-30">
          <div className="flex justify-center border-b-2 border-zinc-100 px-5 py-5">
            <div className="inline-flex items-center justify-center gap-2 rounded-full border-2 border-zinc-200 bg-white px-4 py-1.5 text-xs/none font-semibold uppercase text-zinc-900">
              <Circle
                className="size-2.5 fill-primary text-primary"
                aria-hidden="true"
              />
              <span>Documentation</span>
            </div>
          </div>

          {nav}

          <div className="relative z-10 mx-auto w-full lg:w-178">
            <div className="border-x-2 border-zinc-100 bg-white px-6 py-10 md:px-10 md:py-15">
              {children}
            </div>
          </div>
        </div>
      </section>
    </SiteShell>
  );
}
