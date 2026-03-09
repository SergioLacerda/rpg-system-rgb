# Penetração

Penetração representa a capacidade de uma arma **superar a proteção de armaduras** no Sistema RGB.

Armas com maior poder balístico ou impacto físico possuem maior capacidade de atravessar proteções.

## Cálculo de Penetração

A penetração reduz temporariamente o valor de armadura do alvo durante o cálculo de dano.

```text
Armadura Efetiva = Armadura − Penetração
```

Exemplo:

```text
Armadura = 4
Penetração da arma = 2

Armadura efetiva = 2
```

Se o valor de armadura efetiva chegar a **0 ou menos**, a armadura não reduz o dano.

## Tabela de Penetração de Armas

```text
Arma              Penetração
----------------  --------Pistola           1
SMG               2
Carabina          2
Rifle de Assalto  3
Rifle de Precisão 4
Metralhadora      3
```

## Interação com o Modelo de Dano

A penetração é aplicada **antes da redução de dano pela armadura**.

Fluxo completo de resolução de dano:

```text
Dano da Arma
↓
Penetração
↓
Armadura Efetiva
↓
Absorção por Escudo
↓
Dano Restante → Personagem
```

## Interação com o Sistema RGB

Penetração está relacionada principalmente ao impacto da arma, mas interage indiretamente com os vetores RGB.

```text
R → aumenta dano em armas corpo a corpo
G → permite evitar ataques através de mobilidade
B → fornece proteção adicional através de escudos
```

Armas com maior penetração são mais eficazes contra alvos fortemente protegidos.

## Filosofia de Design

O sistema de penetração no RGB segue três princípios:

- **simplicidade** — cálculo direto e fácil de aplicar
- **equilíbrio** — armas diferentes possuem capacidades distintas contra armaduras
- **compatibilidade** — integração direta com o modelo de dano RGB

Isso permite diferenciar armas sem adicionar complexidade excessiva ao combate.

## Veja Também

- [Armas de Fogo](../categories/firearms.md)
- [Modelo de Dano](../../combat/damage_model.md)
- [Armadura](../../equipment/armor.md)

← [Voltar para Armas de Fogo](../README.md)
