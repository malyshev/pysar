import { HomeClients } from "@/components/home-clients";
import { HomeFeatures } from "@/components/home-features";
import { HomeHero } from "@/components/home-hero";
import { HomeIntegrations } from "@/components/home-integrations";
import { HomeSteps } from "@/components/home-steps";
import { JsonLd } from "@/components/json-ld";
import { SiteShell } from "@/components/site-shell";
import { buildHomeJsonLd } from "@/lib/json-ld";
import { buildHomeMetadata } from "@/lib/metadata";

export const metadata = buildHomeMetadata();

export default function HomePage() {
  return (
    <SiteShell framed={false}>
      <JsonLd data={buildHomeJsonLd()} />
      <HomeHero />
      <HomeClients />
      <HomeFeatures />
      <HomeSteps />
      <HomeIntegrations />
    </SiteShell>
  );
}
