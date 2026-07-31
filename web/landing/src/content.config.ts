import { defineCollection, z } from 'astro:content';
import { glob } from 'astro/loaders';

const library = defineCollection({
  loader: glob({ pattern: '**/*.md', base: './src/content/library' }),
  schema: z.object({
    title: z.string(),
    section: z.string(),
    vector: z.enum(['R', 'G', 'B']).optional(),
    status: z.enum(['canonical', 'draft']).default('draft'),
    version: z.string().default('0.1.0'),
    lang: z.enum(['pt-BR', 'en']).default('pt-BR')
  })
});

export const collections = { library };
