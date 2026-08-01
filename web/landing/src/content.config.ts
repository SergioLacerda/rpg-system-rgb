import { defineCollection } from "astro:content";
import { glob } from "astro/loaders";
import { z } from "astro/zod";

const library = defineCollection({
  loader: glob({ pattern: "**/*.md", base: "./src/content/library" }),
  schema: z.object({
    title: z.string(),
    section: z.string(),
    vector: z.enum(["R", "G", "B"]).optional(),
    status: z.enum(["canonical", "draft"]).default("draft"),
    version: z.string().default("0.1.0"),
    lang: z.enum(["pt-BR", "en"]).default("pt-BR"),
  }),
});

export const collections = { library };
