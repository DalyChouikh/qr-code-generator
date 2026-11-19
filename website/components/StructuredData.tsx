import { releasesUrl, repoUrl, siteDescription, siteName, siteUrl } from "@/lib/site";

const FEATURES = [
  "QR codes for links, WiFi networks, contact cards (vCard), email, SMS and plain text",
  "PNG and SVG output from 64 to 4096 pixels",
  "Custom foreground and background colors",
];

const structuredData = {
  "@context": "https://schema.org",
  "@graph": [
    {
      "@type": "WebApplication",
      name: `${siteName} QR Code Generator`,
      url: `${siteUrl}/`,
      description: siteDescription,
      applicationCategory: "UtilitiesApplication",
      operatingSystem: "Any",
      browserRequirements: "Requires a browser with WebAssembly support",
      isAccessibleForFree: true,
      offers: { "@type": "Offer", price: "0", priceCurrency: "USD" },
      featureList: FEATURES,
    },
    {
      "@type": "SoftwareApplication",
      name: siteName,
      description: "Interactive terminal QR code generator with a step-by-step wizard, history and self-update.",
      applicationCategory: "DeveloperApplication",
      operatingSystem: "macOS, Linux, Windows",
      downloadUrl: releasesUrl,
      codeRepository: repoUrl,
      license: `${repoUrl}/blob/main/LICENSE`,
      isAccessibleForFree: true,
      offers: { "@type": "Offer", price: "0", priceCurrency: "USD" },
      featureList: FEATURES,
    },
  ],
};

// Escaping "<" keeps the JSON from closing the script tag early.
const json = JSON.stringify(structuredData).replace(/</g, "\\u003c");

export function StructuredData() {
  return <script type="application/ld+json" dangerouslySetInnerHTML={{ __html: json }} />;
}
