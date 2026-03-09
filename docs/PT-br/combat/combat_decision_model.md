# Modelo de Decisão de Combate

O Sistema RGB incentiva os jogadores a tomarem decisões táticas durante o combate.

Cada vetor RGB representa uma abordagem estratégica diferente.

## Triângulo Tático RGB

As decisões de combate podem ser visualizadas usando o triângulo RGB.

```text
        G
   mobilidade / posicionamento

R ---------------- B
poder / dano      escudo / defesa
```

Cada vetor representa uma estratégia de combate diferente.

## Estratégias de Combate

### R – Estratégia Ofensiva

Foco em dano direto e combate agressivo.

Ações típicas:

- ataques poderosos
- armas pesadas
- quebrar defesas inimigas

Vantagens:

- alto dano
- eliminação rápida de inimigos

Fraqueza:

- menor sobrevivência

### G – Estratégia de Mobilidade

Foco em movimento e posicionamento.

Ações típicas:

- flanquear inimigos
- evitar ataques
- controlar distância

Vantagens:

- mais difícil de atingir
- melhor posicionamento tático

Fraqueza:

- menor dano direto

### B – Estratégia Defensiva

Foco em sobrevivência e proteção.

Ações típicas:

- uso de escudos
- habilidades defensivas
- absorção de ataques

Vantagens:

- alta durabilidade
- grande capacidade defensiva

Fraqueza:

- resolução de combate mais lenta

## Interações Táticas

O sistema RGB permite múltiplos estilos de combate.

```text
| Estratégia | Foco do Vetor |
|------------|---------------|
Lutador agressivo | Alto R |
Escaramuçador móvel | Alto G |
Guardião defensivo | Alto B |
```

Construções híbridas também são possíveis.

Exemplos:

```text
| Build | Estilo |
|------|--------|
R + G | atacante rápido |
R + B | guerreiro pesado |
G + B | defensor evasivo |
```

## Fluxo de Decisão de Combate

Durante o combate um jogador geralmente escolhe entre três opções táticas:

```text
Atacar
Mover
Defender
```

Essas ações se relacionam naturalmente com os vetores RGB.

```text
Atacar → R
Mover → G
Defender → B
```

## Interação com o Modelo de Dano

As decisões de combate interagem diretamente com o sistema de dano.

```text
Ataque
↓
Dano da Arma
↓
Penetração
↓
Armadura
↓
Escudo
↓
Personagem
```

O jogador deve escolher entre:

- aumentar o dano
- evitar dano
- absorver dano

## Filosofia de Design

O Sistema RGB cria um modelo tático simples:

```text
R → causar dano
G → evitar dano
B → absorver dano
```

Essas três estratégias formam o núcleo das decisões táticas em combate no RGB.
