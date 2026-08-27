import { describe, expect, it } from "vitest";
import { dict, locales } from "./dictionary";

describe("landing dictionary", () => {
  it("declares the supported landing locales in routing order", () => {
    expect(locales).toEqual(["pt-br", "en"]);
  });

  it("keeps core Portuguese copy available for the default home", () => {
    expect(dict["pt-br"].hero.title).toBe(
      "Todo conflito se resolve em três vetores.",
    );
    expect(dict["pt-br"].footer.tag).toBe("código aberto");
    expect(dict["pt-br"].install.title).toBe("Rodar o projeto localmente");
  });

  it("keeps the English landing copy unchanged", () => {
    expect(dict.en.hero.title).toBe(
      "Every conflict resolves along three vectors.",
    );
    expect(dict.en.footer.tag).toBe("open source");
    expect(dict.en.install.title).toBe("Run the project locally");
  });

  it("provides one visible vector card per RGB vector for every locale", () => {
    for (const locale of locales) {
      expect(dict[locale].vectors.map((vector) => vector.code)).toEqual([
        "R",
        "G",
        "B",
      ]);
      expect(dict[locale].vectors).toHaveLength(3);
    }
  });

  it("does not present contract-only skills as installed runtime products", () => {
    for (const locale of locales) {
      const skills = dict[locale].skills;

      expect(skills.specialist.name).toBe("Specialist");
      expect(skills.specialist.badge).not.toMatch(/installed|instalada/i);
      expect(skills.maker.badge).not.toMatch(/installed|instalada/i);
      // Maker has no runtime yet and its copy says so explicitly; Specialist
      // is packaged (see the `packaged` test below) and its badge/tags carry
      // the "contract defined" signal instead of a runtime disclaimer.
      expect(skills.maker.desc).toMatch(/runtime/i);
    }
  });

  it("only marks the specialist skill as packaged (a real .zip exists)", () => {
    for (const locale of locales) {
      const skills = dict[locale].skills;

      expect(skills.specialist.packaged).toBe(true);
      expect(skills.maker.packaged).toBe(false);
    }
  });
});
