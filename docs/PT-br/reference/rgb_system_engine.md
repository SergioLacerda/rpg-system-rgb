# Engine do Sistema RGB

O **Engine do Sistema RGB** descreve a arquitetura completa do sistema de RPG RGB.
Ele conecta todos os componentes do sistema em uma estrutura em camadas semelhante à de uma engine de jogos.

Esse modelo ajuda a explicar como o sistema RGB permanece:

- modular
- escalável
- fácil de balancear
- adaptável a diversos cenários

## Arquitetura do Sistema RGB

O sistema RGB pode ser entendido como uma engine em camadas:

```text
           Cenário da Campanha
                │
                ▼
        RGB Ability Engine
                │
                ▼
          Camada de Gameplay
                │
                ▼
           Sistema de Combate
                │
                ▼
            Motor de Dano
                │
                ▼
           Dados de Equipamento
                │
                ▼
            Vetores RGB
```

Cada camada se apoia na camada abaixo.

## Descrição das Camadas

## Vetores RGB (Camada Base)

Todo o sistema é construído sobre três vetores fundamentais.

```text
        G
   mobilidade / posicionamento

R ---------------- B
poder / dano        escudo / defesa
```

Os vetores definem as mecânicas principais:

- **R (Vermelho)** → poder ofensivo
- **G (Verde)** → mobilidade e reação
- **B (Azul)** → defesa e energia

Todos os outros sistemas derivam desses três valores.

## Dados de Equipamento

Essa camada define os valores numéricos utilizados pelo sistema de combate.

Exemplos:

- dano de armas
- redução de armadura
- valores de penetração
- capacidade de escudos

Documentos relacionados:

- firearms.md
- melee.md
- armor.md
- shields.md

## Motor de Dano

O motor de dano determina como ataques afetam os personagens.

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

Documentos relacionados:

- damage_model.md
- rgb_damage_interaction_model.md

## Sistema de Combate

A camada de combate define a interação entre personagens.

Exemplos:

- resolução de ataques
- movimento
- posicionamento
- reações defensivas

Documentos relacionados:

- movement.md
- attack_and_defense.md

## Camada de Gameplay

Define como os jogadores interagem com o sistema durante o jogo.

Loop típico de jogo:

```text
Criação de Personagem
↓
Exploração
↓
Combate
↓
Recuperação
```

Documentos relacionados:

- character_creation.md
- gameplay_loop.md
- combat_decision_model.md

## RGB Ability Engine

Essa camada introduz habilidades, poderes e mecânicas especiais.

As habilidades são construídas usando os vetores RGB.

```text
Vetor
Custo
Efeito
Duração
Limite
```

Documentos relacionados:

- skills_and_abilities.md
- rgb_ability_engine.md

## Cenário da Campanha

A camada superior define o mundo narrativo.

Exemplos:

- campanhas modernas
- mundos de fantasia
- ficção científica
- universos com superpoderes

Exemplo de separação:

```text
Sistema RGB → engine de regras
Aurora → cenário construído sobre a engine
```

## Fluxo Completo do Sistema

A interação completa entre as camadas:

```text
Vetores
 ↓
Dados de Equipamento
 ↓
Motor de Dano
 ↓
Sistema de Combate
 ↓
Regras de Gameplay
 ↓
Engine de Habilidades
 ↓
Cenário da Campanha
```

## Vantagens de Design

O RGB System Engine oferece diversas vantagens.

## Design Modular

Novos subsistemas podem ser adicionados sem reescrever as regras centrais.

Exemplos:

- sistemas de magia
- veículos
- cibernética
- superpoderes

## Balanceamento Mais Simples

A maior parte do balanceamento ocorre em:

```text
Motor de Dano
Dados de Equipamento
```

O restante do sistema permanece estável.

## Documentação Clara

Cada camada responde a uma pergunta específica.

```text
Camada            Pergunta
--------------------------------------------Vetores           O que define um personagem?
Equipamentos      Quais números existem?
Motor de Dano     Como o dano funciona?
Combate           Como personagens interagem?
Gameplay          Como jogadores interagem?
Habilidades       Quais poderes existem?
Cenário           Em que mundo a história ocorre?
```

## Observação Final

O Sistema RGB se comporta menos como um conjunto tradicional de regras de RPG
e mais como uma **engine modular de RPG**.

Graças a essa arquitetura, o mesmo sistema pode suportar diferentes gêneros
mantendo simplicidade e consistência.

← [Voltar para README](README.md)
