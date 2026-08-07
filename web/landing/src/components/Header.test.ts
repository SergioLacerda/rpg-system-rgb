import { experimental_AstroContainer as AstroContainer } from "astro/container";
import { describe, expect, it } from "vitest";
import Header from "./Header.astro";

describe("Header", () => {
  it("renders Portuguese language controls with the default home unprefixed", async () => {
    const container = await AstroContainer.create();
    const html = await container.renderToString(Header, {
      props: { lang: "pt-br" },
    });

    expect(html).toContain('href="/rpg-system-rgb/" class="active"');
    expect(html).toContain('href="/rpg-system-rgb/en/"');
    expect(html).toContain("PT-BR");
    expect(html).toContain("EN");
    expect(html).not.toContain("Sobre");
  });

  it("renders English language controls without changing the English route", async () => {
    const container = await AstroContainer.create();
    const html = await container.renderToString(Header, {
      props: { lang: "en" },
    });

    expect(html).toContain('href="/rpg-system-rgb/en/" class="active"');
    expect(html).toContain('href="/rpg-system-rgb/"');
    expect(html).toContain("Install");
    expect(html).not.toContain("About");
  });
});
