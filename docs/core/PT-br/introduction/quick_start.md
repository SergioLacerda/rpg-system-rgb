# Sistema RGB – Início Rápido

Este guia de início rápido explica as ideias centrais do **Sistema RGB** e como começar a jogar em poucos minutos.

## 1. Os Vetores RGB

Os personagens no Sistema RGB são definidos por três vetores:

```text
| Vetor | Nome | Função |
|------|------|--------|
| R | Vermelho | pressão, impacto, disrupção |
| G | Verde | relação, movimento, reação |
| B | Azul | preservação, escudos, estabilização |
```

Esses três valores definem como um personagem atua em combate e interação.

## 2. Criando um Personagem

Um personagem inicial normalmente recebe:

```text
7 pontos
```

Os jogadores distribuem esses pontos entre **R, G e B**.

Exemplo:

```text
R = 3
G = 2
B = 2
```

## 3. Durabilidade do Personagem

A durabilidade é calculada usando fórmulas simples.

```text
Vida = 4 + R + B
Escudo = B × 3
```

Exemplo:

```text
R = 3
B = 2
Vida = 9
Escudo = 6
```

## 4. Estrutura Básica de Turno

Durante o combate um personagem normalmente tem:

```text
Movimento (G × 2 metros, livre)
+
1 Ação
+
1 Ação Menor (opcional)
```

Exemplos:

Ações:

- atacar
- defender
- usar habilidade

Ações menores:

- recarregar arma
- mudar posição levemente
- interagir com o ambiente

## 5. Decisões de Combate

O combate gira em torno de três escolhas táticas:

```text
Atacar
Mover
Defender
```

Essas ações correspondem diretamente aos vetores RGB:

```text
Pressionar   → R
Reposicionar → G
Sustentar    → B
```

Os jogadores naturalmente enfatizam diferentes estratégias dependendo da distribuição de seus vetores.

## 6. Resolução de Dano

Quando um ataque acerta, o dano segue esta sequência:

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

Armadura reduz o dano **por ataque**, enquanto escudos absorvem **dano acumulado**.

## 7. Exemplo de Situação de Combate

Um personagem com:

```text
R = 3
G = 2
B = 2
```

enfrenta um inimigo.

O jogador deve escolher entre:

- aumentar o dano causado
- reposicionar-se
- reforçar a defesa

Cada decisão corresponde a um dos vetores RGB.

## 8. Loop de Jogo

Uma sessão típica de RGB segue esta estrutura:

```text
Criação de Personagem
↓
Exploração
↓
Encontro de Combate
↓
Recuperação
↓
Progressão da História
```

Esse ciclo se repete conforme a história avança.

## Saiba Mais

Para uma compreensão mais profunda veja:

- Regras Básicas → [Regras Básicas](../core/README.md)
- Sistema de Combate → [Sistema de Combate](../combat/README.md)
- Equipamentos → [Equipamentos](../equipment/README.md)
- Armas → [Armas](../weapons/README.md)
- Loop de Jogo → [Loop de Jogo](../reference/gameplay_loop.md)
