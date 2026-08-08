import { homePageSeo, siteConfig, siteLogo } from "@/lib/site";
import { getLogoUrl, getSiteUrl } from "@/lib/metadata";

export function buildHomeJsonLd() {
  const siteUrl = getSiteUrl();
  const orgId = `${siteUrl}/#organization`;
  const websiteId = `${siteUrl}/#website`;

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
    ],
  };
}
