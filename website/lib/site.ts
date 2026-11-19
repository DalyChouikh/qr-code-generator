export const repo = "DalyChouikh/qr-code-generator";
export const repoUrl = `https://github.com/${repo}`;
export const releasesUrl = `${repoUrl}/releases/latest`;
export const rawUrl = `https://raw.githubusercontent.com/${repo}/main`;

// Absolute URLs are required for canonical links, the sitemap and social cards; set this when deploying.
export const siteUrl = (process.env.NEXT_PUBLIC_SITE_URL ?? "http://localhost:3000").replace(/\/$/, "");

export const siteName = "qrgen";
export const siteTitle = "QRgen - QR Code Generator";
export const siteDescription =
  "Free, open-source QR code generator. Make QR codes for links, WiFi networks, contact cards (vCard), email and SMS, then download them as PNG or SVG. Runs in your browser, nothing is uploaded, and the same generator ships as a terminal app.";
