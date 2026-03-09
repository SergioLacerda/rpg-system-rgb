# Escudos

Escudos são **equipamentos defensivos adicionais** que aumentam a proteção de um personagem no Sistema RGB.

Eles funcionam em conjunto com **Armaduras** e o **Modelo de Dano**, criando camadas de defesa.

Existem dois tipos principais de escudos:

- Escudos Físicos
- Escudos Energéticos / Mágicos

## Escudos Físicos

Escudos físicos representam equipamentos defensivos carregados por um personagem.

Eles fornecem proteção adicional, mas podem impor penalidades de mobilidade.

Escudos físicos podem ser **destruídos, derrubados ou perdidos** dependendo da decisão do Mestre.

```text
Tipo de Escudo   Proteção   Penalidade de Mobilidade
---------------  ---------- ------------------------
Escudo Leve      +1         0
Escudo Médio     +2         −1 G
Escudo Pesado    +3         −2 G
```

### Características

- Adiciona proteção à defesa
- Pode afetar mobilidade (G)
- Sujeito a dano ou perda de equipamento

Escudos físicos são mais comuns em:

- cenários medievais
- ambientes de combate tático
- combate em curta distância

## Escudos Energéticos ou Mágicos

Escudos energéticos ou mágicos representam **campos de energia protetores**, tecnologia avançada ou barreiras mágicas.

Diferente da armadura, esses escudos **absorvem dano em vez de reduzi-lo**.

Características:

- não regeneram durante o combate
- possuem valor fixo durante um confronto
- retornam ao valor máximo após o término do combate

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

Esse sistema de defesa em camadas mantém o combate simples enquanto preserva profundidade tática.

## Interação com os Vetores RGB

Escudos interagem com os três vetores RGB de maneiras diferentes.

```text
R → determina resistência e vida
G → pode ser reduzido por escudos pesados
B → determina a capacidade do escudo energético
```

Isso mantém a filosofia de combate do RGB:

```text
R → causar dano
G → evitar dano
B → absorver dano
```

## Filosofia de Design

Escudos no sistema RGB seguem três princípios:

- **defesa em camadas** – armadura reduz dano, escudos absorvem dano
- **trocas táticas** – escudos mais pesados reduzem mobilidade
- **design modular** – escudos se adaptam a diferentes cenários

Exemplos:

- campanhas medievais → escudos físicos
- ficção científica → escudos energéticos
- cenários de fantasia → escudos mágicos

## Veja Também

- [Armadura](armor.md)
- [Modelo de Dano](../combat/damage_model.md)
- [Ataque e Defesa](../combat/attack_and_defense.md)

← [Voltar para README](README.md)
