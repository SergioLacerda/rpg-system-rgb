# Armas Corpo a Corpo

Armas corpo a corpo representam **armas utilizadas em combate direto** no Sistema RGB.

Diferente de armas de fogo, armas corpo a corpo podem **utilizar o vetor R
(Vermelho) do personagem como parte da Fonte de Impacto** quando a arma ou
procedimento permitir.

Isso representa força física, impacto e habilidade marcial do usuário.

## Fonte de Impacto

A Fonte de Impacto de uma arma corpo a corpo normalmente é calculada somando o
valor de **R do personagem** ao bônus da arma.

```text
Fonte de Impacto = R + bônus da arma
```

Exemplo:

```text
R = 3
Espada Longa = +3

Fonte de Impacto = 6
```

## Tabela de Armas Corpo a Corpo

```text
Arma            Bônus de Dano
--------------  -----------Faca            +1
Adaga           +1
Espada Curta    +2
Espada Longa    +3
Machado         +3
```

## Interação com o Sistema RGB

Armas corpo a corpo estão diretamente ligadas ao vetor **R (Vermelho)**.

```text
R → determina a pressão e pode contribuir para a Fonte de Impacto
G → influencia mobilidade, timing e posicionamento
B → preserva continuidade por bloqueios, escudos e estabilização
```

Personagens com valores altos de **R** tendem a ser mais eficazes em combate corpo a corpo.

## Interação com o Modelo de Dano

Quando um ataque corpo a corpo acerta um alvo, o fluxo de dano segue o modelo padrão do RGB.

```text
Fonte de Impacto (R + arma)
↓
Penetração (se aplicável)
↓
Redução por Armadura
↓
Absorção por Escudo
↓
Dano Restante → Personagem
```

## Uso Tático

Armas corpo a corpo são particularmente eficazes em:

- combate em curta distância
- espaços confinados
- situações onde armas de fogo não podem ser utilizadas

Elas também permitem combate silencioso ou ataques surpresa.

## Filosofia de Design

Armas corpo a corpo no RGB seguem três princípios:

- **simplicidade** — cálculo direto baseado no vetor R
- **equilíbrio** — dano proporcional à força do personagem
- **compatibilidade** — integração direta com o modelo de dano RGB

## Veja Também

- [Combate](../../combat/attack_and_defense.md)
- [Modelo de Dano](../../combat/damage_model.md)
- [Armas de Fogo](../categories/firearms.md)

← [Voltar para README](../README.md)
