# Armadura

Armaduras representam **proteção física** usada para reduzir o dano recebido no Sistema RGB.

A armadura interage com o **Modelo de Dano** e também pode se combinar com **escudos energéticos** caso o vetor Azul (B) esteja ativo.

## Cálculo de Armadura

A proteção da armadura pode ser entendida como:

```text
Defesa Total = Armadura (proteção física) + Escudo energético (ou mágico, dependendo do cenário) (pontos B)
```

**Observação:**  
O Escudo energético (ou mágico, dependendo do cenário) é aplicado apenas se o **Fator Azul (B)** estiver ativo no cenário da campanha.

Exemplo:

```text
Armadura = 4
Escudo = B × 3
```

A armadura reduz o dano por ataque, enquanto o escudo absorve dano acumulado.

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

Usuários de armadura pesada dependem mais de **R (força)** e **B (defesa)** do que de mobilidade.

## Interação da Armadura com o Sistema RGB

A armadura interage com os vetores RGB da seguinte forma:

```text
R → determina resistência física e vida
G → afetado pelas penalidades de mobilidade da armadura
B → pode fornecer proteção adicional através de escudos
```

Essa relação mantém o equilíbrio tático do RGB:

```text
R → causar dano
G → evitar dano
B → absorver dano
```

## Interação com o Motor de Dano

A armadura reduz o dano **após a aplicação da penetração da arma**.

Fluxo simplificado de dano:

```text
Dano da Arma
↓
Penetração
↓
Redução por Armadura
↓
Absorção pelo Escudo
↓
Dano Restante → Personagem
```

Para as regras completas veja:

- damage_model.md

## Filosofia de Design

Armaduras no sistema RGB seguem três princípios:

- **valores simples** – fáceis de calcular durante o combate
- **trocas táticas** – maior proteção reduz mobilidade
- **compatibilidade com escudos** – permite camadas de defesa

Isso garante que armaduras permaneçam relevantes sem tornar personagens excessivamente resistentes.

## Veja Também

- [Escudos](shields.md)
- [Modelo de Dano](../combat/damage_model.md)

← [Voltar para README](README.md)
