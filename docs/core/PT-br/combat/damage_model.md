# Modelo de Dano

O Sistema RGB usa um **modelo de dano em camadas** que combina fontes de impacto,
penetração, armaduras, escudos e habilidades especiais. Esse modelo mantém o
combate simples enquanto permite decisões táticas.

Este documento explica como o dano flui pelo sistema.

## Fluxo de Resolução de Dano

O dano no Sistema RGB segue uma sequência clara:

```text
Verificação de Acerto ou Contato
     ↓
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

Esse fluxo mostra como diferentes elementos do sistema interagem durante o
combate.

## Resolução Passo a Passo

### 1. Verificação de Ataque

O atacante tenta acertar o defensor.

As regras de combate determinam se o ataque é bem-sucedido.

Veja:

- [Ataque e Defesa](attack_and_defense.md)

Se o ataque falhar, nenhum dano é aplicado.

### 2. Fonte de Impacto

Se o ataque for bem-sucedido, identifique a **Fonte de Impacto**.

Impacto pode vir de uma arma, atributo, habilidade, procedimento ou exceção
explícita. O padrão é:

- impacto corpo a corpo usa o valor da arma mais o `R` do atacante quando a arma
  ou procedimento permitir;
- impacto de arma de fogo usa o valor da arma por padrão;
- impacto explosivo usa o perfil do explosivo;
- impacto de habilidade usa o efeito declarado pela habilidade.

Exemplos não podem sobrescrever a Fonte de Impacto declarada.

Veja:

- [Armas de Fogo](../weapons/categories/firearms.md)
- [Armas Corpo a Corpo](../weapons/categories/melee.md)
- [Explosivos](../weapons/categories/explosives.md)

### 3. Penetração

Aplique penetração à armadura antes de a armadura reduzir dano.

```text
Armadura Efetiva = Armadura - Penetração
```

Se a armadura efetiva for zero ou menor, a armadura não reduz aquele acerto.

Veja:

- [Penetração](../weapons/mechanics/penetration.md)

### 4. Redução por Armadura

Armadura fornece **proteção física**.

A armadura reduz dano recebido por acerto depois da penetração.

Categorias típicas de armadura incluem:

- Armadura Leve
- Armadura Média
- Armadura Pesada

Veja:

- [Armadura](../equipment/armor.md)

### 5. Absorção por Escudo

Alguns personagens podem possuir escudos.

Escudos absorvem dano restante **depois da redução por armadura** e antes que o
dano alcance a vida ou gere consequência de estado.

Escudos de energia normalmente são calculados como:

```text
Escudo = B × 3
```

Veja:

- [Escudos](../equipment/shields.md)

### 6. Dano Restante

Depois que armadura e escudos são resolvidos, o dano restante é aplicado ao
personagem.

Esse dano afeta a vida do personagem ou cria uma consequência de estado
declarada.

## Habilidades Especiais

Algumas habilidades podem modificar como o dano é aplicado.

Exemplos:

- técnicas de absorção
- escudos de energia
- técnicas marciais defensivas

Veja:

- [Habilidades](../core/skills_and_abilities.md)

## Visão Geral da Interação de Dano

O modelo de dano RGB pode ser resumido como:

```text
Fonte de Impacto
      ↓
Penetração
      ↓
Redução por Armadura
      ↓
Absorção por Escudo
      ↓
Dano Restante
```

Esse modelo em camadas mantém o combate fácil de entender enquanto permite que
equipamentos e habilidades diferentes interajam de modo significativo.

## Modelo de Interação RGB

O sistema também pode ser visualizado como uma interação entre os vetores RGB:

```text
        G
   mobilidade / reação

R - B
poder / dano    escudo / energia
```

Cada vetor representa uma estratégia defensiva ou ofensiva diferente no combate.

- **R (Vermelho)** muda a fonte de pressão por força, impacto ou interrupção.
- **G (Verde)** muda a relação com a pressão por movimento, timing e
  posicionamento.
- **B (Azul)** preserva continuidade por bloqueios, escudos, estabilização e
  resistência.

Juntos, eles criam um sistema tático equilibrado.

← [Voltar para README](README.md)
