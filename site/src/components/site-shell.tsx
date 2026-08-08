import { SiteFooter } from "@/components/site-footer";
import { SiteHeader } from "@/components/site-header";

export function SiteShell({
  children,
  framed = true,
}: {
  children: React.ReactNode;
  /** Docs keep framed main; home uses template section containers. */
  framed?: boolean;
}) {
  return (
    <div className="flex min-h-full flex-1 flex-col">
      <SiteHeader />
      {framed ? (
        <main className="grow px-4 pt-24 md:px-5 lg:px-7.5">
          <div className="site-container flex-1 px-5 py-12 md:px-8 md:py-16">
            {children}
          </div>
        </main>
      ) : (
        <main className="grow">{children}</main>
      )}
      <SiteFooter />
    </div>
  );
}
