import type { Locale } from "../i18n/dictionary";

const testLandingBase = "/rpg-system-rgb/";
const configuredBase =
  import.meta.env.MODE === "test" && import.meta.env.BASE_URL === "/"
    ? testLandingBase
    : (import.meta.env.BASE_URL ?? "/");

export function normalizeLandingBase(base = configuredBase): string {
  if (base === "/") {
    return "";
  }

  return base.endsWith("/") ? base.slice(0, -1) : base;
}

export function landingPath(path: string): string {
  const base = normalizeLandingBase();
  const normalizedPath = path.startsWith("/") ? path : `/${path}`;

  if (normalizedPath === "/") {
    return `${base}/`;
  }

  return `${base}${normalizedPath}`;
}

export function pathWithoutLandingBase(pathname: string): string {
  const base = normalizeLandingBase();

  if (!base) {
    return pathname;
  }

  if (pathname === base) {
    return "/";
  }

  if (pathname.startsWith(`${base}/`)) {
    return pathname.slice(base.length);
  }

  return pathname;
}

export function localizedLandingPath(lang: Locale, rest = ""): string {
  const normalizedRest = rest.replace(/^\/+|\/+$/g, "");

  if (lang === "pt-br" && normalizedRest === "") {
    return landingPath("/");
  }

  const suffix = normalizedRest ? `/${normalizedRest}` : "/";
  return landingPath(`/${lang}${suffix}`);
}

export function localizedRestPath(pathname: string): string {
  const currentPath = pathWithoutLandingBase(pathname);
  const segments = currentPath.split("/").filter(Boolean);

  if (segments[0] === "pt-br" || segments[0] === "en") {
    return segments.slice(1).join("/");
  }

  return segments.join("/");
}
