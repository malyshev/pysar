import {
  agentInstallGuide,
  homePageSeo,
  installCommand,
  siteConfig,
  siteLogo,
} from "@/lib/site";
import { getLogoUrl, getSiteUrl } from "@/lib/metadata";

export function buildHomeJsonLd() {
  const siteUrl = getSiteUrl();
  const orgId = `${siteUrl}/#organization`;
  const websiteId = `${siteUrl}/#website`;
  const softwareId = `${siteUrl}/#software`;
  const howtoId = `${siteUrl}/#howto-agent-install`;
  const installDocs = `${siteUrl}${agentInstallGuide.docsInstallPath}`;
  const initDocs = `${siteUrl}${agentInstallGuide.docsInitPath}`;

  return {
    "@context": "https://schema.org",
    "@graph": [
      {
        "@type": "Organization",
        "@id": orgId,
        name: siteConfig.name,
        url: siteUrl,
        logo: {
          "@type": "ImageObject",
          url: getLogoUrl(),
          width: siteLogo.width,
          height: siteLogo.height,
        },
      },
      {
        "@type": "WebSite",
        "@id": websiteId,
        name: siteConfig.name,
        url: siteUrl,
        description: homePageSeo.description,
        publisher: { "@id": orgId },
      },
      {
        "@type": "SoftwareApplication",
        "@id": softwareId,
        name: siteConfig.name,
        applicationCategory: "Productivity",
        operatingSystem: "macOS, Linux, Windows",
        url: siteUrl,
        downloadUrl: `${siteUrl}/install.sh`,
        installUrl: installDocs,
        description: agentInstallGuide.metaSummary,
        publisher: { "@id": orgId },
        offers: {
          "@type": "Offer",
          price: "0",
          priceCurrency: "USD",
        },
      },
      {
        "@type": "HowTo",
        "@id": howtoId,
        name: "Install Pysar when a user pastes getpysar.com and says install",
        description: agentInstallGuide.metaSummary,
        url: `${siteUrl}/#for-ai-agents`,
        step: [
          {
            "@type": "HowToStep",
            position: 1,
            name: "Install the pysar binary",
            text: `Run: ${installCommand} Then confirm with pysar --version. Do not invent brew, npm, pip, or go install alternatives.`,
            url: installDocs,
          },
          {
            "@type": "HowToStep",
            position: 2,
            name: "Init for the host agent",
            text: agentInstallGuide.hosts
              .map((h) => `${h.name}: ${h.initCommand}`)
              .join(". "),
            url: initDocs,
          },
          {
            "@type": "HowToStep",
            position: 3,
            name: "Continue in the host",
            text: "Open the project in that host and continue with /ps-onboard or /ps as the author directs.",
          },
        ],
      },
    ],
  };
}
