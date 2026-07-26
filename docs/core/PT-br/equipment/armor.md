# Armadura

Armaduras representam **proteção física** usada para reduzir dano recebido no
Sistema RGB.

A armadura interage com o **Modelo de Dano** e pode coexistir com escudos, mas
ela não se soma a escudos para formar uma defesa total genérica. Armadura é uma
camada de redução por acerto; escudo é uma camada posterior de absorção.

## Proteção de Armadura

A proteção da armadura é aplicada como redução por acerto depois da penetração.

```text
Armadura Efetiva = Armadura - Penetração
Dano Após Armadura = Fonte de Impacto - Armadura Efetiva
```

Se a Armadura Efetiva for zero ou menor, a armadura não reduz aquele acerto.

Exemplo:

```text
Fonte de Impacto = 7
Armadura = 4
Penetração = 1
Armadura Efetiva = 3
Dano após armadura = 4
```

A absorção por escudo, se existir, acontece depois dessa etapa.

## Tipos de Armadura

```text
Tipo      Proteção   Penalidade de Mobilidade
--------  ---------- ------------------------
Leve      2          −1 G
Média     4          −2 G
Pesada    6          −3 G
```

### Armadura Leve

- Proteção mínima
- Baixa penalidade de mobilidade
- Adequada para personagens ágeis

### Armadura Média

- Proteção equilibrada
- Redução moderada de mobilidade

### Armadura Pesada

- Alta proteção física
- Penalidade significativa de mobilidade

Usuários de armadura pesada dependem mais de procedimentos de bloqueio,
posicionamento e preservação do que de evasão pura.

## Interação da Armadura com o Sistema RGB

A armadura interage com os vetores RGB da seguinte forma:

```text
R → define ou contribui para a Fonte de Impacto quando a ação permite
G → é afetado pelas penalidades de mobilidade da armadura
B → preserva continuidade por bloqueios, escudos e estabilização
```

Essa relação mantém o equilíbrio tático do RGB:

```text
R → mudar a fonte de pressão
G → mudar a relação com a pressão
B → preservar continuidade sob pressão
```

## Interação com o Motor de Dano

A armadura reduz dano **após a aplicação da penetração**.

Fluxo simplificado de dano:

```text
Fonte de Impacto
↓
Penetração
↓
Redução por Armadura
↓
Absorção por Escudo
↓
Dano Restante → Personagem
```

Para as regras completas veja:

- [Modelo de Dano](../combat/damage_model.md)

## Filosofia de Design

Armaduras no Sistema RGB seguem três princípios:

- **valores simples** — fáceis de calcular durante o combate
- **trocas táticas** — maior proteção reduz mobilidade
- **compatibilidade com escudos** — permite camadas defensivas sem colapsar
  procedimentos diferentes

Isso garante que armaduras permaneçam relevantes sem tornar personagens
excessivamente resistentes.

## Veja Também

- [Escudos](shields.md)
- [Modelo de Dano](../combat/damage_model.md)

← [Voltar para README](README.md)
