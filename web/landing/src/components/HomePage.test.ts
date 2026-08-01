import { experimental_AstroContainer as AstroContainer } from "astro/container";
import { describe, expect, it } from "vitest";
import HomePage from "./HomePage.astro";

describe("HomePage", () => {
  it("renders the Portuguese home with the default latest PDF link", async () => {
    const container = await AstroContainer.create();
    const html = await container.renderToString(HomePage, {
      props: { lang: "pt-br" },
    });

    expect(html).toContain("Todo conflito se resolve em três vetores.");
    expect(html).toContain("Regras demais, jogo de menos.");
    expect(html).toContain(
      "/rpg-system-rgb/downloads/rgb-system-core-v2-latest-pt-br.pdf",
    );
  });

  it("renders the English home copy and PDF link", async () => {
    const container = await AstroContainer.create();
    const html = await container.renderToString(HomePage, {
      props: { lang: "en" },
    });

    expect(html).toContain("Every conflict resolves along three vectors.");
    expect(html).toContain("Too many rules, not enough game.");
    expect(html).toContain(
      "/rpg-system-rgb/downloads/rgb-system-core-v2-latest-en.pdf",
    );
  });
});
