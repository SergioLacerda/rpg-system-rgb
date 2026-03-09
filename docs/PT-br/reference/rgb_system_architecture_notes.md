# Arquitetura do Sistema RGB – Observação de Design em Camadas

Esta nota documenta uma propriedade estrutural interessante que aparece quando a documentação do RGB
é analisada como um todo.

Embora o Sistema RGB tenha sido projetado como um sistema de RPG de mesa simples, a documentação
revela naturalmente uma **arquitetura em camadas** semelhante à encontrada em engines de jogos.

Essa estrutura melhora a clareza, modularidade e capacidade de expansão a longo prazo.

## Arquitetura em Camadas do Sistema RGB

Quando a documentação é organizada por função, o sistema forma as seguintes camadas:

```text
Camada do Sistema
↓
Camada de Gameplay
↓
Camada de Combate
↓
Camada de Dano
↓
Camada de Equipamentos
```

Cada camada responde a uma pergunta diferente de design.

## 1. Camada do Sistema

Documentos:

- system_overview.md
- rgb_damage_interaction_model.md

Objetivo:

Esta camada explica a **filosofia central e a estrutura do sistema**.

Ela define os vetores RGB:

- **R (Vermelho)** – dano e poder físico
- **G (Verde)** – mobilidade e reação
- **B (Azul)** – escudos e defesa energética

Também mostra como os módulos do sistema se conectam.

Equivalente em engines de jogo:

Design do sistema central.

## 2. Camada de Gameplay

Documentos:

- gameplay_loop.md
- combat_decision_model.md

Objetivo:

Esta camada explica **como os jogadores interagem com o sistema durante o jogo**.

Fluxo típico:

```text
Criação de Personagem
↓
Exploração
↓
Combate
↓
Recuperação
```

Durante o combate o jogador normalmente escolhe entre:

```text
Atacar
Mover
Defender
```

Que se relacionam naturalmente com:

```text
Atacar → R
Mover → G
Defender → B
```

Equivalente em engines de jogo:

Estrutura de gameplay.

## 3. Camada de Combate

Documentos:

- attack_and_defense.md
- movement.md

Objetivo:

Define como os personagens interagem durante encontros.

Isso inclui:

- regras de movimento
- resolução de ataques
- reações defensivas
- posicionamento

Equivalente em engines de jogo:

Sistema de combate.

## 4. Camada de Dano

Documentos:

- damage_model.md
- rgb_damage_interaction_model.md

Objetivo:

Responsável pela lógica interna do cálculo de dano.

O motor de dano RGB segue esta estrutura:

```text
Arma
↓
Penetração
↓
Armadura
↓
Escudo
↓
Personagem
```

Essa camada permite que o combate permaneça tático enquanto mantém os cálculos simples.

Equivalente em engines de jogo:

Motor de dano.

## 5. Camada de Equipamentos

Documentos:

- armor.md
- shields.md
- firearms.md
- melee.md
- explosives.md

Objetivo:

Define os **valores de dados utilizados pelos sistemas de combate e dano**.

Exemplos:

- dano de armas
- redução de armadura
- valores de penetração
- capacidade de escudos

Equivalente em engines de jogo:

Camada de dados.

## Arquitetura Completa do Sistema RGB

Quando combinada, a documentação do RGB forma a seguinte estrutura:

```text
Design do Sistema
   │
   ├─ Visão Geral do Sistema
   ├─ Modelo de Interação RGB
   │
Gameplay
   │
   ├─ Loop de Jogo
   ├─ Modelo de Decisão de Combate
   │
Combate
   │
   ├─ Movimento
   ├─ Ataque e Defesa
   │
Motor de Dano
   │
   ├─ Modelo de Dano
   │
Dados de Equipamentos
   │
   ├─ Armaduras
   ├─ Escudos
   ├─ Armas
```

## Consequências de Design

Essa arquitetura em camadas produz diversas vantagens.

### Modularidade

Novos subsistemas podem ser adicionados sem modificar as regras centrais.

Exemplos:

- sistemas de magia
- cibernética
- veículos
- superpoderes

### Balanceamento Mais Simples

A maior parte do balanceamento ocorre em:

```text
Camada de Dano
Camada de Equipamentos
```

O restante do sistema permanece estável.

### Documentação Clara

Cada camada responde a uma pergunta diferente:

```text
| Camada | Pergunta |
|--------|----------|
Sistema | Como o sistema é estruturado |
Gameplay | Como os jogadores interagem com o jogo |
Combate | Como os personagens interagem |
Dano | Como o dano é calculado |
Equipamentos | Quais valores numéricos existem |
```

## Observação Final

O Sistema RGB se comporta menos como um conjunto tradicional de regras de RPG
e mais como uma **engine reutilizável de sistema de RPG**.

Isso torna o sistema particularmente adequado para:

- campanhas modernas
- mundos de fantasia
- cenários de ficção científica
- universos com superpoderes

Nesse modelo:

```text
Sistema RGB → engine de regras
Cenário (ex.: Aurora) → mundo construído sobre a engine
```

Essa separação permite que o Sistema RGB permaneça genérico,
enquanto diferentes cenários fornecem o contexto narrativo.

← [Voltar para README](README.md)
