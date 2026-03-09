# Sistema RGB — Índice da Documentação (PT-BR)

Este índice central organiza toda a documentação do **Sistema RGB** em português.

O objetivo é facilitar a navegação entre os módulos principais do sistema.

## RGB System

Sistema de RPG mantendo regras simples e combate tático.

O **Sistema RGB** utiliza três vetores:

```text
Vetor   Nome       Função
- - --
R       Vermelho   ataques, potência, teste de força, acerto, dano aplicado.
G       Verde      agilidade, mobililidade, esquiva, velocidade de reação.
B       Azul       escudo magico, proteção especial, intelecto, sabedoria, energia, resistência
```

Foi projetado com quatro princípios:

- simplicidade estrutural
- profundidade tática
- modularidade
- escalabilidade

Ele permite adaptação para:

- campanhas modernas
- fantasia
- ficção científica
- superpoderes

## Introdução

Documentos que apresentam os conceitos fundamentais do sistema.

- [Visão Geral do Sistema](introduction/system_overview.md)
- [Início Rápido](introduction/quick_start.md)
- [Regras em Uma Página](introduction/rgb_one_page_rules.md)

## Sistema Base

Define os fundamentos do sistema RGB.

- [Atributos](core/attributes.md)
- [Criação de Personagem](core/character_creation.md)
- [Habilidades](core/skills_and_abilities.md)

## Combate

Regras que governam encontros de combate.

- [Movimento](combat/movement.md)
- [Ataque e Defesa](combat/attack_and_defense.md)
- [Modelo de Dano](combat/damage_model.md)
- [Modelo de Decisão de Combate](combat/combat_decision_model.md)

## Equipamentos

Equipamentos que afetam sobrevivência e tática.

- [Armaduras](equipment/armor.md)
- [Escudos](equipment/shields.md)
- [Equipamentos Gerais](equipment/gear.md)

## Armas

Tipos de armas e regras relacionadas.

- [Armas de Fogo](weapons/firearms.md)
- [Armas Corpo a Corpo](weapons/melee.md)
- [Explosivos](weapons/explosives.md)
- [Penetração](weapons/penetration.md)
- [Alcance de Armas](weapons/weapon_ranges.md)

## Referência

Materiais auxiliares e exemplos práticos.

- [Ficha de Personagem](reference/character_sheet.md)
- [Exemplo de Combate](reference/combat_example.md)
- [Exemplo de Jogo](reference/gameplay_example.md)
- [Loop de Jogo](reference/gameplay_loop.md)
- [Engine do Sistema RGB](reference/rgb_system_engine.md)

## Notas de Design

Documentos que explicam decisões de design e arquitetura do sistema.

- [Modelo de Interação RGB](reference/rgb_damage_interaction_model.md)
- [Arquitetura do Sistema RGB](reference/rgb_system_architecture_notes.md)
- [Notas Extras](reference/extra.md)

## Estrutura da Documentação

```text
PT-br/
├── introduction/
├── core/
├── combat/
├── equipment/
├── weapons/
└── reference/
```

Cada pasta possui um **README.md próprio** que funciona como índice local da seção.
