# Exemplo de Combate

Este exemplo demonstra como o combate é resolvido usando o **Sistema RGB**.

Dois personagens entram em combate e resolvem ações usando as regras principais:

- vetores RGB
- movimento
- ataque e defesa
- armadura e escudos
- cálculo de dano

## Personagens

## Personagem A – Combatente de Assalto

```text
R = 3
G = 2
B = 1
```

Valores derivados:

```text
Vida = 4 + (R × 2) = 10
Escudo = B × 3 = 3
Movimento = G × 2 = 4 metros
```

Equipamento:

- Armadura Média (Proteção 4)
- Rifle

## Personagem B – Operador Móvel

```text
R = 2
G = 3
B = 2
```

Valores derivados:

```text
Vida = 8
Escudo = 6
Movimento = 6 metros
```

Equipamento:

- Armadura Leve (Proteção 2)
- SMG

## Exemplo de Turno de Combate

## Passo 1 — Posicionamento

Os personagens se movem de acordo com:

```text
Movimento = G × 2 metros
```

O Personagem B se move primeiro devido à maior mobilidade.

## Passo 2 — Ataque

O Personagem A dispara o rifle.

Exemplo de dano da arma:

```text
Dano da Arma = 6
Penetração = 2
```

## Passo 3 — Interação com Armadura

A armadura reduz o dano recebido após a penetração.

Exemplo:

```text
Armadura do Alvo = 2
Penetração = 2
Armadura Efetiva = 0
```

A armadura é ignorada.

## Passo 4 — Absorção por Escudo

O escudo absorve dano antes da vida.

```text
Escudo = 6
Dano = 6
```

Resultado:

```text
Escudo Restante = 0
Dano à Vida = 0
```

## Segundo Ataque

O Personagem B ataca com uma SMG.

```text
Dano da Arma = 4
Penetração = 1
```

Armadura do Personagem A:

```text
Armadura = 4
Armadura Efetiva = 3
```

Cálculo de dano:

```text
Dano = 4 - 3 = 1
```

O Personagem A recebe:

```text
Dano à Vida = 1
Vida Restante = 9
```

## Resumo do Fluxo de Dano

O sistema RGB resolve dano usando uma estrutura em camadas.

```text
Dano da Arma
↓
Penetração
↓
Redução por Armadura
↓
Absorção por Escudo
↓
Dano Restante → Personagem
```

## Interpretação Tática

O exemplo mostra o papel dos vetores RGB:

```text
R → determina poder ofensivo
G → determina mobilidade e posicionamento
B → determina capacidade de escudo
```

Cada vetor sustenta uma estratégia de combate diferente.

## Objetivo de Design

A resolução de combate no sistema RGB foi projetada para ser:

- rápida de calcular
- taticamente significativa
- consistente entre cenários

Como as mesmas fórmulas se aplicam a todos os personagens, o sistema permanece fácil de balancear.

## Veja Também

- [Ataque e Defesa](../combat/attack_and_defense.md)
- [Movimento](../combat/movement.md)
- [Modelo de Dano](../combat/damage_model.md)

← [Voltar para README](README.md)
