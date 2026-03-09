# Sistema RGB – Início Rápido

Este guia de início rápido explica as ideias centrais do **Sistema RGB** e como começar a jogar em poucos minutos.

## 1. Os Vetores RGB

Os personagens no Sistema RGB são definidos por três vetores:

```text
| Vetor | Nome | Função |
|------|------|--------|
| R | Vermelho | dano, força, poder ofensivo |
| G | Verde | agilidade, movimento, reação |
| B | Azul | escudos, energia, defesa especial |
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
Vida = 4 + (R × 2)
Escudo = B × 3
```

Exemplo:

```text
R = 3 → Vida = 10
B = 2 → Escudo = 6
```

## 4. Estrutura Básica de Turno

Durante o combate um personagem normalmente pode realizar:

```text
1 ação principal
+
1 ajuste menor (opcional)
```

Exemplos:

Ações principais:

- atacar
- mover
- defender
- usar habilidade

Ajustes menores:

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
Atacar → R
Mover → G
Defender → B
```

Os jogadores naturalmente enfatizam diferentes estratégias dependendo da distribuição de seus vetores.

## 6. Resolução de Dano

Quando um ataque acerta, o dano segue esta sequência:

```text
Teste de Ataque
↓
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
