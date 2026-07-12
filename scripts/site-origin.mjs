// Site origin for absolute OGP URLs (SNS crawlers require absolute og:image / og:url).
// Override with SITE_ORIGIN when deploying to a custom domain.
export const SITE_ORIGIN = (
  process.env.SITE_ORIGIN ||
  process.env.EBISHOWCASE_ORIGIN ||
  "https://kumagi.github.io/EbiShowcase"
).replace(/\/$/, "");

export function absoluteURL(pathname) {
  const path = String(pathname || "").replace(/^\/+/, "");
  return path ? `${SITE_ORIGIN}/${path}` : `${SITE_ORIGIN}/`;
}
