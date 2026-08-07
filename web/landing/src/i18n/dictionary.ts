export const locales = ["pt-br", "en"] as const;
export type Locale = (typeof locales)[number];

export const dict = {
  "pt-br": {
    nav: {
      install: "Instalação",
      library: "Library",
      skills: "Skills",
    },
    a11y: {
      sectionNav: {
        hero: "Início",
        problem: "O PROBLEMA",
        system: "RGB System",
        library: "LIBRARY",
        forDevelopers: "PARA DESENVOLVEDORES",
      },
      scrollHint: "Avançar",
      scrollHintPrev: "Voltar",
    },
    hero: {
      kicker: "UM SISTEMA DE RPG DE MESA · CÓDIGO ABERTO",
      title: "Todo conflito se resolve em três vetores.",
      sub: "Pressão, Relação e Preservação substituem dezenas de subsistemas — regras enxutas, decisões com peso real.",
      ctaLibrary: "Ler o material — Library",
      ctaPdf: "Baixar PDF",
    },
    vectors: [
      {
        code: "R",
        label: "VETOR R · PRESSÃO",
        title: "Ação sob risco",
        desc: "Regras de conflito, dano e consequência imediata — o quanto o mundo empurra de volta.",
      },
      {
        code: "G",
        label: "VETOR G · RELAÇÃO",
        title: "Vínculos e influência",
        desc: "Como personagens afetam o outro — negociação, lealdade, reputação com peso mecânico.",
      },
      {
        code: "B",
        label: "VETOR B · PRESERVAÇÃO",
        title: "Recursos e desgaste",
        desc: "Gestão do que se perde com o tempo: energia, materiais, integridade do personagem.",
      },
    ],
    problem: {
      kicker: "O PROBLEMA",
      title: "Regras demais, jogo de menos.",
      items: [
        "Sistemas com dezenas de subsistemas para uma única decisão de combate.",
        "Notas de campanha espalhadas, sem estrutura ou fonte única.",
        "Fichas e regras que divergem entre mesas e entre sessões.",
        "Material de referência difícil de consultar no meio do jogo.",
      ],
    },
    system: {
      kicker: "RGB SYSTEM",
      title: "Três vetores. Regras leves. Decisões táticas.",
      desc: "Toda ação do jogo passa por R, G ou B — sem subsistemas paralelos para aprender.",
    },
    forDevelopers: {
      kicker: "PARA DESENVOLVEDORES",
      title: "O sistema é dados abertos, não só texto.",
      desc: "Regras, procedimentos e relações entre eles são publicados como um índice semântico versionado, compilado em um bundle JSON consumível por ferramentas.",
      items: [
        { label: "Bundle JSON", code: "generated/bundle/rgb.bundle.json" },
        { label: "Compilar o bundle", code: "make bundle" },
        { label: "Validar a documentação", code: "make validate" },
      ],
    },
    libraryTeaser: {
      kicker: "LIBRARY",
      title: "O material completo, sempre atualizado.",
      desc: "Cada capítulo do sistema publicado como página navegável — com status de versão, exemplos e link direto para o Markdown fonte. Sem cadastro, sem paywall.",
      cta: "Abrir Library →",
      previewTitle: "Resolução de Ataque",
      previewDesc:
        "O acerto é determinado pela margem entre o vetor ofensivo e a resistência do alvo...",
    },
    library: {
      kicker: "LIBRARY",
      title: "Regras navegáveis, com fonte aberta.",
      desc: "A Library será a superfície de leitura do RGB System: capítulos versionados, status editorial e links diretos para os arquivos Markdown canônicos.",
      empty: "Selecione um capítulo.",
      pdf: "Baixar PDF do capítulo",
      source: "Ver em Markdown",
      sections: [
        {
          vector: "R",
          title: "Pressão",
          desc: "Conflito, dano, risco imediato e consequências de ação.",
        },
        {
          vector: "G",
          title: "Relação",
          desc: "Influência, vínculos, reputação e negociação com peso mecânico.",
        },
        {
          vector: "B",
          title: "Preservação",
          desc: "Recursos, desgaste, proteção e continuidade do personagem.",
        },
      ],
    },
    install: {
      kicker: "INSTALAÇÃO · PASSO 1",
      title: "Rodar o projeto localmente",
      cloneLabel: "CLONE + INSTALL",
      devLabel: "DEV SERVER",
      copy: "COPY",
    },
    skills: {
      kicker: "SKILLS",
      title: "Agentes para consultar e produzir material RGB.",
      desc: "As skills ficam fora da landing, mas a página apresenta o papel operacional de cada uma dentro do ecossistema.",
      specialist: {
        name: "Specialist",
        status: "contrato definido",
        desc: "Contrato para futura consulta de regras do RGB System. O runtime ainda não está implementado; a skill deve citar IDs semânticos e caminhos Markdown quando existir.",
        tags: [
          "runtime não implementado",
          "cita IDs semânticos",
          "cita fonte .md",
        ],
      },
      maker: {
        name: "Maker",
        status: "contrato definido",
        desc: "Contrato para futura estruturação de notas, imagens e aventuras em documentos. O runtime ainda não está implementado e não cria cânone automaticamente.",
        tags: [
          "runtime não implementado",
          "separa fatos e sugestões",
          "não cria cânone",
        ],
      },
    },
    footer: { tag: "código aberto" },
  },
  en: {
    nav: {
      install: "Install",
      library: "Library",
      skills: "Skills",
    },
    a11y: {
      sectionNav: {
        hero: "Home",
        problem: "THE PROBLEM",
        system: "RGB System",
        library: "LIBRARY",
        forDevelopers: "FOR DEVELOPERS",
      },
      scrollHint: "Next section",
      scrollHintPrev: "Previous section",
    },
    hero: {
      kicker: "A TABLETOP RPG SYSTEM · OPEN SOURCE",
      title: "Every conflict resolves along three vectors.",
      sub: "Pressure, Relation and Preservation replace dozens of subsystems — lean rules, decisions that matter.",
      ctaLibrary: "Read the rules — Library",
      ctaPdf: "Download PDF",
    },
    vectors: [
      {
        code: "R",
        label: "VECTOR R · PRESSURE",
        title: "Action under risk",
        desc: "Conflict, damage and immediate consequence rules — how hard the world pushes back.",
      },
      {
        code: "G",
        label: "VECTOR G · RELATION",
        title: "Bonds and influence",
        desc: "How characters affect one another — negotiation, loyalty, reputation with mechanical weight.",
      },
      {
        code: "B",
        label: "VECTOR B · PRESERVATION",
        title: "Resources and wear",
        desc: "Managing what erodes over time: energy, supplies, a character\u2019s integrity.",
      },
    ],
    problem: {
      kicker: "THE PROBLEM",
      title: "Too many rules, not enough game.",
      items: [
        "Systems with dozens of subsystems for a single combat decision.",
        "Campaign notes scattered around, with no single structured source.",
        "Sheets and rules that drift apart between tables and between sessions.",
        "Reference material that is hard to look up mid-session.",
      ],
    },
    system: {
      kicker: "RGB SYSTEM",
      title: "Three vectors. Lean rules. Tactical decisions.",
      desc: "Every action in the game routes through R, G, or B — no parallel subsystems to learn.",
    },
    forDevelopers: {
      kicker: "FOR DEVELOPERS",
      title: "The system is open data, not just prose.",
      desc: "Rules, procedures, and the relationships between them are published as a versioned semantic index, compiled into a consumable JSON bundle.",
      items: [
        { label: "JSON bundle", code: "generated/bundle/rgb.bundle.json" },
        { label: "Build the bundle", code: "make bundle" },
        { label: "Validate the docs", code: "make validate" },
      ],
    },
    libraryTeaser: {
      kicker: "LIBRARY",
      title: "The full material, always up to date.",
      desc: "Every chapter published as a browsable page — version status, examples, and a direct link to the source Markdown. No signup, no paywall.",
      cta: "Open Library →",
      previewTitle: "Attack Resolution",
      previewDesc:
        "A hit is determined by the margin between the offensive vector and the target resistance...",
    },
    library: {
      kicker: "LIBRARY",
      title: "Browsable rules, open at the source.",
      desc: "The Library is the reading surface for RGB System: versioned chapters, editorial status, and direct links back to the canonical Markdown files.",
      empty: "Select a chapter.",
      pdf: "Download chapter PDF",
      source: "View Markdown",
      sections: [
        {
          vector: "R",
          title: "Pressure",
          desc: "Conflict, damage, immediate risk, and action consequences.",
        },
        {
          vector: "G",
          title: "Relation",
          desc: "Influence, bonds, reputation, and negotiation with mechanical weight.",
        },
        {
          vector: "B",
          title: "Preservation",
          desc: "Resources, wear, protection, and character continuity.",
        },
      ],
    },
    install: {
      kicker: "INSTALL · STEP 1",
      title: "Run the project locally",
      cloneLabel: "CLONE + INSTALL",
      devLabel: "DEV SERVER",
      copy: "COPY",
    },
    skills: {
      kicker: "SKILLS",
      title: "Agents for consulting and producing RGB material.",
      desc: "Skills live outside the landing page, but this page presents each operational role in the ecosystem.",
      specialist: {
        name: "Specialist",
        status: "contract defined",
        desc: "Contract for future RGB System rule consultation. Runtime behavior is not implemented yet; the skill must cite semantic IDs and Markdown paths when it exists.",
        tags: [
          "runtime not implemented",
          "cites semantic IDs",
          "cites source .md",
        ],
      },
      maker: {
        name: "Maker",
        status: "contract defined",
        desc: "Contract for future structuring of notes, images, and adventures into documents. Runtime behavior is not implemented yet and does not create canon automatically.",
        tags: [
          "runtime not implemented",
          "separates facts and suggestions",
          "does not create canon",
        ],
      },
    },
    footer: { tag: "open source" },
  },
} as const;
