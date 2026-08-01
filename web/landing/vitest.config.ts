/// <reference types="vitest" />

import { getViteConfig } from "astro/config";

export default getViteConfig({
  base: "/rpg-system-rgb/",
  test: {
    environment: "node",
    include: ["src/**/*.test.ts"],
    coverage: {
      provider: "v8",
      reporter: ["text", "lcov"],
      include: [
        "src/components/**/*.astro",
        "src/i18n/**/*.ts",
        "src/lib/**/*.ts",
      ],
      exclude: [
        "src/**/*.test.ts",
        "src/content.config.ts",
        "src/env.d.ts",
        "src/styles/**",
      ],
      thresholds: {
        statements: 90,
        branches: 90,
        functions: 90,
        lines: 90,
      },
    },
  },
} as Parameters<typeof getViteConfig>[0] & {
  test: Record<string, unknown>;
});
