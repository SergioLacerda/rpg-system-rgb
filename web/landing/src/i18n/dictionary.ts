export const locales = ['pt-br', 'en'] as const;
export type Locale = typeof locales[number];

export const dict = {
  'pt-br': {
    nav: { about: 'Sobre', install: 'Instalação', library: 'Library', skills: 'Skills' },
    hero: {
      kicker: 'UM SISTEMA DE RPG DE MESA · CÓDIGO ABERTO',
      title: 'Todo conflito se resolve em três vetores.',
      sub: 'Pressão, Relação e Preservação substituem dezenas de subsistemas — regras enxutas, decisões com peso real.',
      ctaLibrary: 'Ler o material — Library',
      ctaPdf: 'Baixar PDF (v0.2)'
    },
    vectors: [
      { code: 'R', label: 'VETOR R · PRESSÃO', title: 'Ação sob risco', desc: 'Regras de conflito, dano e consequência imediata — o quanto o mundo empurra de volta.' },
      { code: 'G', label: 'VETOR G · RELAÇÃO', title: 'Vínculos e influência', desc: 'Como personagens afetam o outro — negociação, lealdade, reputação com peso mecânico.' },
      { code: 'B', label: 'VETOR B · PRESERVAÇÃO', title: 'Recursos e desgaste', desc: 'Gestão do que se perde com o tempo: energia, materiais, integridade do personagem.' }
    ],
    libraryTeaser: {
      kicker: 'LIBRARY',
      title: 'O material completo, sempre atualizado.',
      desc: 'Cada capítulo do sistema publicado como página navegável — com status de versão, exemplos e link direto para o Markdown fonte. Sem cadastro, sem paywall.',
      cta: 'Abrir Library →',
      previewTitle: 'Resolução de Ataque',
      previewDesc: 'O acerto é determinado pela margem entre o vetor ofensivo e a resistência do alvo...'
    },
    library: {
      kicker: 'LIBRARY',
      title: 'Regras navegáveis, com fonte aberta.',
      desc: 'A Library será a superfície de leitura do RGB System: capítulos versionados, status editorial e links diretos para os arquivos Markdown canônicos.',
      empty: 'Selecione um capítulo.',
      pdf: 'Baixar PDF do capítulo',
      source: 'Ver em Markdown',
      sections: [
        { vector: 'R', title: 'Pressão', desc: 'Conflito, dano, risco imediato e consequências de ação.' },
        { vector: 'G', title: 'Relação', desc: 'Influência, vínculos, reputação e negociação com peso mecânico.' },
        { vector: 'B', title: 'Preservação', desc: 'Recursos, desgaste, proteção e continuidade do personagem.' }
      ]
    },
    install: {
      kicker: 'INSTALAÇÃO · PASSO 1',
      title: 'Rodar o projeto localmente',
      cloneLabel: 'CLONE + INSTALL',
      devLabel: 'DEV SERVER',
      copy: 'COPY'
    },
    skills: {
      kicker: 'SKILLS',
      title: 'Agentes para consultar e produzir material RGB.',
      desc: 'As skills ficam fora da landing, mas a página apresenta o papel operacional de cada uma dentro do ecossistema.',
      pathfinder: { name: 'Pathfinder', status: 'instalada', desc: 'Especialista na documentação do RGB System. Responde dúvidas sobre regras, indica o capítulo exato da Library e resolve ambiguidades entre versões.', tags: ['busca em toda a Library', 'explica vetores R·G·B', 'cita fonte .md'] },
      maker: { name: 'Maker', status: 'não instalada', desc: 'Converte histórias e aventuras em rascunho para documentos estruturados — padrão RGB System por default, exportável para outros sistemas.', tags: ['texto bruto → .md estruturado', 'mapeia p/ vetores R·G·B', 'export multi-sistema'] }
    },
    footer: { tag: 'código aberto' }
  },
  en: {
    nav: { about: 'About', install: 'Install', library: 'Library', skills: 'Skills' },
    hero: {
      kicker: 'A TABLETOP RPG SYSTEM · OPEN SOURCE',
      title: 'Every conflict resolves along three vectors.',
      sub: 'Pressure, Relation and Preservation replace dozens of subsystems — lean rules, decisions that matter.',
      ctaLibrary: 'Read the rules — Library',
      ctaPdf: 'Download PDF (v0.2)'
    },
    vectors: [
      { code: 'R', label: 'VECTOR R · PRESSURE', title: 'Action under risk', desc: 'Conflict, damage and immediate consequence rules — how hard the world pushes back.' },
      { code: 'G', label: 'VECTOR G · RELATION', title: 'Bonds and influence', desc: 'How characters affect one another — negotiation, loyalty, reputation with mechanical weight.' },
      { code: 'B', label: 'VECTOR B · PRESERVATION', title: 'Resources and wear', desc: 'Managing what erodes over time: energy, supplies, a character\u2019s integrity.' }
    ],
    libraryTeaser: {
      kicker: 'LIBRARY',
      title: 'The full material, always up to date.',
      desc: 'Every chapter published as a browsable page — version status, examples, and a direct link to the source Markdown. No signup, no paywall.',
      cta: 'Open Library →',
      previewTitle: 'Attack Resolution',
      previewDesc: 'A hit is determined by the margin between the offensive vector and the target resistance...'
    },
    library: {
      kicker: 'LIBRARY',
      title: 'Browsable rules, open at the source.',
      desc: 'The Library is the reading surface for RGB System: versioned chapters, editorial status, and direct links back to the canonical Markdown files.',
      empty: 'Select a chapter.',
      pdf: 'Download chapter PDF',
      source: 'View Markdown',
      sections: [
        { vector: 'R', title: 'Pressure', desc: 'Conflict, damage, immediate risk, and action consequences.' },
        { vector: 'G', title: 'Relation', desc: 'Influence, bonds, reputation, and negotiation with mechanical weight.' },
        { vector: 'B', title: 'Preservation', desc: 'Resources, wear, protection, and character continuity.' }
      ]
    },
    install: {
      kicker: 'INSTALL · STEP 1',
      title: 'Run the project locally',
      cloneLabel: 'CLONE + INSTALL',
      devLabel: 'DEV SERVER',
      copy: 'COPY'
    },
    skills: {
      kicker: 'SKILLS',
      title: 'Agents for consulting and producing RGB material.',
      desc: 'Skills live outside the landing page, but this page presents each operational role in the ecosystem.',
      pathfinder: { name: 'Pathfinder', status: 'installed', desc: 'Documentation specialist for the RGB System. Answers rules questions, points to the exact Library chapter, and resolves ambiguity between versions.', tags: ['searches the whole Library', 'explains R·G·B vectors', 'cites the source .md'] },
      maker: { name: 'Maker', status: 'not installed', desc: 'Converts raw stories and adventures into structured documents — RGB System by default, exportable to other systems.', tags: ['raw text → structured .md', 'maps to R·G·B vectors', 'multi-system export'] }
    },
    footer: { tag: 'open source' }
  }
} as const;
