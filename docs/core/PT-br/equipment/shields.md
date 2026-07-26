# Escudos

Escudos são **equipamentos ou recursos defensivos adicionais** que ajudam um
personagem a preservar continuidade sob pressão no Sistema RGB.

Eles funcionam com **Armaduras** e o **Modelo de Dano**, criando camadas de
defesa. Escudos não se somam a armadura para formar uma defesa total genérica;
eles absorvem dano restante depois da redução por armadura, salvo regra
específica em contrário.

Existem dois tipos principais de escudos:

- Escudos Físicos
- Escudos Energéticos / Mágicos

## Escudos Físicos

Escudos físicos representam equipamentos defensivos carregados por um
personagem.

Eles permitem bloquear, proteger aliados ou absorver dano restante, mas podem
impor penalidades de mobilidade.

Escudos físicos podem ser **destruídos, derrubados ou perdidos** dependendo da
decisão do Mestre.

```text
Tipo de Escudo   Proteção   Penalidade de Mobilidade
---------------  ---------- ------------------------
Escudo Leve      +1         0
Escudo Médio     +2         −1 G
Escudo Pesado    +3         −2 G
```

### Características

- Apoia procedimentos de bloqueio ou proteção
- Pode afetar mobilidade (G)
- Sujeito a dano ou perda de equipamento

Escudos físicos são mais comuns em:

- cenários medievais
- ambientes de combate tático
- combate em curta distância

## Escudos Energéticos ou Mágicos

Escudos energéticos ou mágicos representam **campos de energia protetores**,
tecnologia avançada ou barreiras mágicas.

Diferente da armadura, esses escudos **absorvem dano em vez de reduzi-lo por
acerto**.

Características:

- não regeneram durante o combate, salvo regra específica
- possuem valor fixo durante um confronto
- retornam ao valor máximo após o término do combate quando o cenário permite

O valor do escudo é determinado pelo vetor Azul.

```text
Escudo = B × 3
```

Exemplo:

```text
B = 3
Escudo = 9
```

Escudos energéticos representam a **camada de absorção** no sistema de dano RGB.

## Interação com o Modelo de Dano RGB

Os escudos são aplicados **após a redução de dano pela armadura**.

Ordem de resolução do dano:

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

Esse sistema de defesa em camadas mantém o combate simples enquanto preserva
profundidade tática.

## Interação com os Vetores RGB

Escudos interagem com os três vetores RGB de maneiras diferentes.

```text
R → pode definir a fonte de pressão que o escudo tenta conter
G → pode ser reduzido por escudos pesados
B → determina a capacidade de escudos energéticos e preservação
```

Isso mantém a filosofia de combate do RGB:

```text
R → mudar a fonte de pressão
G → mudar a relação com a pressão
B → preservar continuidade sob pressão
```

## Filosofia de Design

Escudos no Sistema RGB seguem três princípios:

- **defesa em camadas** — armadura reduz dano, escudos absorvem dano
- **trocas táticas** — escudos mais pesados reduzem mobilidade
- **design modular** — escudos se adaptam a diferentes cenários

Exemplos:

- campanhas medievais → escudos físicos
- ficção científica → escudos energéticos
- cenários de fantasia → escudos mágicos

## Veja Também

- [Armadura](armor.md)
- [Modelo de Dano](../combat/damage_model.md)
- [Ataque e Defesa](../combat/attack_and_defense.md)

← [Voltar para README](README.md)
