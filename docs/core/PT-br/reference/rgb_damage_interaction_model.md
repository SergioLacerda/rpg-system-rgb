# Sistema RGB – Modelo de Dano e Notas de Interação

Este documento consolida observações de design sobre o modelo de dano do Sistema
RGB, seu equilíbrio matemático e como os vetores RGB interagem com as mecânicas
de combate.

O objetivo é explicar **por que o sistema se comporta de forma consistente e
permanece equilibrado**, mantendo regras simples.

## Modelo de Interação RGB

O Sistema RGB é construído em torno de três vetores fundamentais:

- **R (Vermelho)** — pressão, impacto e disrupção
- **G (Verde)** — relação, mobilidade e reação
- **B (Azul)** — preservação, escudos e estabilização

Esses vetores interagem para criar um combate tático.

```text
        G
   mobilidade / reação

R - B
pressão          preservação
```

Cada vetor influencia um aspecto diferente do jogo:

```text
| Vetor | Papel no Jogo |
|------|----------------|
R | muda a fonte de pressão |
G | muda a relação com a pressão |
B | preserva continuidade sob pressão |
```

## Fluxo de Interação do Sistema

O Sistema RGB conecta atributos do personagem, equipamentos e mecânicas de
combate.

```text
Vetores RGB
↓
Criação de Personagem
↓
Seleção de Equipamentos
↓
Interação de Combate
↓
Resolução de Dano
```

Essa estrutura modular permite que o sistema se adapte a diferentes cenários:

- campanhas modernas
- mundos de fantasia
- cenários de ficção científica
- universos com superpoderes

## Modelo de Dano RGB

O dano no combate segue uma estrutura em camadas.

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

Essa abordagem garante que fontes de impacto, armaduras e escudos interajam de
forma previsível.

## Fórmulas Principais do Sistema

O sistema RGB define durabilidade do personagem com fórmulas simples.

```text
Vida = 4 + R + B
Escudo = B × 3
```

A Fonte de Impacto pode vir de arma, atributo, habilidade, procedimento ou
exceção explícita. Armaduras oferecem **valores fixos de redução por acerto**.

A penetração modifica o valor efetivo da armadura.

## Armadura Efetiva

A penetração interage com a armadura da seguinte forma:

```text
Armadura Efetiva = Armadura − Penetração
```

O dano restante após armadura torna-se:

```text
Dano Restante = Fonte de Impacto − Armadura Efetiva
```

Isso mantém os cálculos simples e previsíveis.

## Escala Natural de Impacto

Se utilizarmos os seguintes intervalos:

```text
| Elemento | Escala Típica |
|----------|---------------|
Fontes de Impacto | 3–10 |
Armadura | 1–6 |
Penetração | 0–4 |
```

Surge uma relação natural de equilíbrio:

```text
Fonte de Impacto ≈ Armadura + Penetração
```

Isso garante que ataques e defesas permaneçam relevantes.

## Exemplo

### Fonte Leve

```text
Fonte de Impacto = 4
Penetração = 1
Armadura = 4
```

```text
Armadura Efetiva = 4 − 1 = 3
Dano Restante = 4 − 3 = 1
```

### Fonte Pesada

```text
Fonte de Impacto = 8
Penetração = 3
Armadura = 6
```

```text
Armadura Efetiva = 6 − 3 = 3
Dano Restante = 8 − 3 = 5
```

## Interação com Escudos

Após a redução pela armadura, os escudos absorvem o dano restante.

```text
Dano Restante − Escudo
```

Escudos seguem a fórmula RGB:

```text
Escudo = B × 3
```

Isso cria duas camadas defensivas:

```text
| Tipo de Defesa | Função |
|----------------|--------|
Armadura | reduz dano por acerto |
Escudo | absorve dano acumulado restante |
```

Essa separação simplifica a resolução do combate.

## Regra Prática de Balanceamento

Ao projetar fontes de impacto e armaduras, a seguinte aproximação mantém o
combate equilibrado:

```text
Armadura ≈ Fonte de Impacto / 2
Penetração ≈ Fonte de Impacto / 3
```

Isso mantém a maioria dos ataques relevantes sem tornar as defesas inúteis.

## Vantagem Inesperada: Combate Contra Múltiplos Inimigos

Uma propriedade interessante do modelo de dano RGB é que ele funciona bem quando
personagens enfrentam **múltiplos oponentes simultaneamente**.

Muitos sistemas de RPG têm dificuldades nesse cenário porque valores defensivos
não escalam bem contra ataques repetidos.

No RGB, o modelo em camadas distribui o dano naturalmente:

```text
Fonte de Impacto → Penetração → Armadura → Escudo → Personagem
```

### Armadura

A armadura reduz dano **em cada ataque individual**, evitando que inimigos fracos
sobrecarreguem rapidamente personagens bem protegidos.

### Escudo

Escudos absorvem **o dano total acumulado restante**, funcionando como um
amortecedor contra vários ataques menores.

## Consequências Táticas

Isso cria uma dinâmica de combate estável.

```text
| Tipo de Inimigo | Resultado |
|-----------------|-----------|
Muitos inimigos fracos | Armadura reduz a maioria dos ataques |
Poucos inimigos fortes | Penetração torna-se mais importante |
Inimigos com energia | Escudos tornam-se críticos |
```

Como a armadura funciona por golpe e os escudos funcionam cumulativamente, o
sistema permanece estável mesmo com muitos atacantes.

## Conclusão de Design

A estrutura de dano do RGB forma uma hierarquia clara de defesa:

```text
Fonte de Impacto
↓
Penetração
↓
Redução por Armadura
↓
Absorção por Escudo
↓
Personagem
```

Essa estrutura:

- mantém os cálculos simples
- suporta combate tático
- escala naturalmente com equipamentos
- lida bem com encontros com múltiplos inimigos
