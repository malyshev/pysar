import type { Metadata } from "next";
import { Geist, Roboto } from "next/font/google";
import { buildSiteMetadataDefaults } from "@/lib/metadata";
import { siteConfig, umamiPixelSrc } from "@/lib/site";
import "./globals.css";

const geist = Geist({
  subsets: ["latin"],
  variable: "--font-geist",
  weight: ["400", "500", "600", "700", "800", "900"],
});

const roboto = Roboto({
  subsets: ["latin"],
  variable: "--font-roboto",
  weight: ["400", "500", "700", "900"],
});

export const metadata: Metadata = {
  title: {
    default: siteConfig.name,
    template: `%s · ${siteConfig.name}`,
  },
  description: siteConfig.description,
  metadataBase: new URL(siteConfig.url),
  ...buildSiteMetadataDefaults(),
};

export default function RootLayout({ children }: LayoutProps<"/">) {
  return (
    <html
      lang="en"
      className={`${geist.variable} ${roboto.variable} h-full antialiased`}
      data-scroll-behavior="smooth"
    >
      <body className="flex min-h-full flex-col font-sans">
        {children}
        {/* Umami Cloud pageview pixel — 1×1 GIF at cloud.umami.is/p/… */}
        {/* eslint-disable-next-line @next/next/no-img-element */}
        <img
          src={umamiPixelSrc}
          alt=""
          width={1}
          height={1}
          className="pointer-events-none fixed size-px opacity-0"
          decoding="async"
        />
      </body>
    </html>
  );
}
