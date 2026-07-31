const configuredBase = import.meta.env.BASE_URL ?? '/';

function normalizedBase(): string {
  if (configuredBase === '/') {
    return '';
  }

  return configuredBase.endsWith('/')
    ? configuredBase.slice(0, -1)
    : configuredBase;
}

export function landingPath(path: string): string {
  const base = normalizedBase();
  const normalizedPath = path.startsWith('/') ? path : `/${path}`;

  if (normalizedPath === '/') {
    return `${base}/`;
  }

  return `${base}${normalizedPath}`;
}

export function pathWithoutLandingBase(pathname: string): string {
  const base = normalizedBase();

  if (!base) {
    return pathname;
  }

  if (pathname === base) {
    return '/';
  }

  if (pathname.startsWith(`${base}/`)) {
    return pathname.slice(base.length);
  }

  return pathname;
}
