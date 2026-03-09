# Armas Corpo a Corpo

Armas corpo a corpo representam **armas utilizadas em combate direto** no Sistema RGB.

Diferente de armas de fogo, armas corpo a corpo **utilizam o vetor R (Vermelho) do personagem como base de dano**.

Isso representa força física, impacto e habilidade marcial do usuário.

## Cálculo de Dano

O dano de uma arma corpo a corpo é calculado somando o valor de **R do personagem** ao bônus da arma.

```text
Dano = R + bônus da arma
```

Exemplo:

```text
R = 3
Espada Longa = +3

Dano total = 6
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
R → determina o dano base
G → influencia mobilidade e posicionamento
B → fornece defesa adicional através de escudos
```

Personagens com valores altos de **R** tendem a ser mais eficazes em combate corpo a corpo.

## Interação com o Modelo de Dano

Quando um ataque corpo a corpo acerta um alvo, o fluxo de dano segue o modelo padrão do RGB.

```text
Dano (R + arma)
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
