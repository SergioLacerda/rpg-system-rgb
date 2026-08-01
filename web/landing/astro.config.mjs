import { defineConfig } from "astro/config";

const base = process.env.ASTRO_BASE ?? "/rpg-system-rgb";

// Override ASTRO_BASE for local preview or alternate deployment targets.
export default defineConfig({
  site: "https://sergiolacerda.github.io",
  base,
  i18n: {
    defaultLocale: "pt-br",
    locales: ["pt-br", "en"],
    routing: { prefixDefaultLocale: true },
  },
});
