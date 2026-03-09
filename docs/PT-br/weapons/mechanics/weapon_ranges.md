# Alcance de Armas

O alcance define **a eficácia de armas à distância em diferentes distâncias de combate** no Sistema RGB.

Essas regras permitem representar perda de precisão e eficiência conforme o alvo se afasta.

## Modificadores de Distância

Ataques com armas de fogo recebem modificadores dependendo da distância ao alvo.

```text
Distância   Modificador
---------   --------Curto       0
Médio       −1 R
Longo       −2 R
```

Esses modificadores representam a dificuldade crescente de acertar alvos distantes.

## Tabela de Alcance (metros)

```text
Arma               Curto   Médio   Longo
-----------------  -----   -----   --Pistola Compacta   10      20      35
Pistola            15      30      50
Revólver           20      40      60
SMG                20      50      80
Rifle de Assalto   40      120     250
Rifle Pesado       60      200     400
Rifle de Precisão  150     500     1000
Metralhadora       60      200     400
```

## Alcance Extremo (Opcional)

Ataques além do alcance longo ainda podem ser tentados, mas são significativamente mais difíceis.

```text
Alcance extremo → −4 R
```

O Mestre decide se o alvo ainda está visível e se o disparo é plausível.

## Interação com o Sistema RGB

O alcance se relaciona principalmente com o vetor **R (Vermelho)**, que representa precisão ofensiva.

```text
R → precisão e eficácia do ataque
G → posicionamento para obter melhor linha de tiro
B → proteção defensiva contra ataques
```

Personagens com maior valor de **R** tendem a lidar melhor com penalidades de distância.

## Interação com o Sistema de Combate

Os modificadores de alcance são aplicados durante o ataque.

Fluxo simplificado:

```text
Ataque
↓
Aplicar modificador de distância
↓
Resolver defesa
↓
Aplicar modelo de dano
```

Isso mantém o sistema consistente com:

- attack_and_defense.md
- damage_model.md

## Filosofia de Design

O sistema de alcance do RGB segue três princípios:

- **simplicidade** — poucos modificadores fáceis de lembrar
- **consistência** — mesmo modelo para todas as armas
- **equilíbrio tático** — armas diferentes possuem alcances distintos

Isso permite que o posicionamento tenha impacto real no combate.

## Veja Também

- [Armas de Fogo](../categories/firearms.md)
- [Penetração](penetration.md)
- [Modelo de Dano](../../combat/damage_model.md)

← [Voltar para Readme](../README.md)
