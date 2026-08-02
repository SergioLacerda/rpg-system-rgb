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
      expect(skills.specialist.status).not.toMatch(/installed|instalada/i);
      expect(skills.maker.status).not.toMatch(/installed|instalada/i);
      expect(skills.specialist.desc).toMatch(/runtime/i);
      expect(skills.maker.desc).toMatch(/runtime/i);
    }
  });
});
