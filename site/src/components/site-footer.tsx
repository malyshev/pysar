import { siteConfig } from "@/lib/site";

export function SiteFooter() {
  return (
    <footer className="border-t border-border/80">
      <div className="mx-auto flex w-full max-w-3xl flex-col gap-2 px-4 py-8 text-sm text-muted-foreground md:px-6">
        <p>{siteConfig.name}</p>
        <p>MIT licensed. Bring your take — it does not post on your behalf.</p>
      </div>
    </footer>
  );
}
