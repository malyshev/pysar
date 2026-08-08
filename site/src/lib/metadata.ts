import type { Metadata } from "next";
import { agentInstallGuide, homePageSeo, siteConfig, siteLogo } from "@/lib/site";

const DEFAULT_OG_IMAGE_PATH = "/og-default.png";

export const DEFAULT_OG_IMAGE_ALT =
  "Pysar — writing help that keeps you in charge";

export function getSiteUrl(): string {
  return siteConfig.url.replace(/\/$/, "");
}

export function getSiteMetadataBase(): URL {
  return new URL(`${getSiteUrl()}/`);
}

function absoluteUrl(path: string): string {
  if (path.startsWith("http://") || path.startsWith("https://")) {
    return path;
  }
  const base = getSiteUrl();
  return path.startsWith("/") ? `${base}${path}` : `${base}/${path}`;
}

export function getDefaultOgImageUrl(): string {
  return absoluteUrl(DEFAULT_OG_IMAGE_PATH);
}

export function getLogoUrl(): string {
  return absoluteUrl(siteLogo.path);
}

export function buildSiteMetadataDefaults(): Pick<
  Metadata,
  "openGraph" | "twitter"
> {
  const image = getDefaultOgImageUrl();

  return {
    openGraph: {
      siteName: siteConfig.name,
      locale: "en_US",
      type: "website",
      images: [
        {
          url: image,
          alt: DEFAULT_OG_IMAGE_ALT,
        },
      ],
    },
    twitter: {
      card: "summary_large_image",
      images: [image],
    },
  };
}

export function buildHomeMetadata(): Metadata {
  const url = getSiteUrl();
  const image = getDefaultOgImageUrl();

  return {
    title: {
      absolute: homePageSeo.title,
    },
    description: homePageSeo.description,
    alternates: {
      canonical: url,
    },
    openGraph: {
      title: homePageSeo.title,
      description: homePageSeo.description,
      url,
      siteName: siteConfig.name,
      locale: "en_US",
      type: "website",
      images: [
        {
          url: image,
          alt: DEFAULT_OG_IMAGE_ALT,
        },
      ],
    },
    twitter: {
      card: "summary_large_image",
      title: homePageSeo.title,
      description: homePageSeo.description,
      images: [image],
    },
    // Machine-facing install carrier co-located with the homepage URL
    // (dec-20260808-homepage-ai-install-signals-9eb5a6a9).
    other: {
      "pysar:agent-install": agentInstallGuide.metaSummary,
      "pysar:agent-trigger": agentInstallGuide.trigger,
    },
  };
}
