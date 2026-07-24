
# Habilidades

Habilidades representam capacidades especiais além das ações normais do sistema RGB.
Habilidades podem ser adquiridas por treino, compra ou seleção na criação do personagem.
O uso dessas habilidades é opcional e depende do cenário definido pelo mestre.
O mestre pode limitar o uso de habilidades por combate, cena ou dia.

## Estrutura de Habilidades RGB

As habilidades no sistema RGB são organizadas em torno dos **três vetores fundamentais**.

```text
        G
   mobilidade / posicionamento

R ---------------- B
poder / dano      escudo / defesa
```

Cada vetor forma naturalmente uma **árvore de habilidades**.

- **R (Vermelho)** → habilidades ofensivas
- **G (Verde)** → habilidades de mobilidade
- **B (Azul)** → habilidades defensivas / de energia

Essa estrutura permite que o sistema escale naturalmente com diferentes cenários.

## Meta Habilidades ou magias

Meta Habilidades se referem a talentos sobrenaturais que podem melhorar talentos existentes
ou criar algo completamente novo, dependendo do universo criado.
O mestre deve pensar na habilidade atrelada ao RGB, bem como a sua progressão a cada nível adquirido.

Exemplo:
Precisão: +1 R em ataques à distância a cada 3 níveis
Estilo samurai: a cada nível o usuário domina uma escola de arte marciais envolvendo espadas

## Vermelho (R)

- Dobrar potência
- Ataque retardatário
- Pressão de ar
- Golpe de Poder
- Quebra de Armadura

Exemplo de progressão:

```text
Nível 3 → +1 R em ataques à distância (Precisão)
Nível 6 → Golpe de poder duplo
Nível 9 → Ataque de impacto em área
```

## Verde (G)

- Reflexos ampliados
- Aceleração momentânea
- Passo fantasma (movimento instantâneo curto)
- Desvio automático
- Corrida sobrenatural
- Salto ampliado
- Movimento silencioso
- Distorção de movimento (ataques erram por velocidade)
- Esquiva sobrenatural

## Azul (B)

- Escudo de energia
- Escudo atordoante
- Escudo inteligente / armadura viva
- Invisibilidade
- Clones
- Absorção de Energia

## Fatores do sistema

### Fator Marcial

Fator Marcial permite técnicas defensivas como bloqueio total (Defesa Completa) e Absorção Especial.
A utilização dessas técnicas pode exigir uma ação, reação ou postura defensiva,
de acordo com a decisão do mestre.

### Fator Azul

Caso o mestre tenha um universo onde magias ou meta-habilidades sejam permitidos.
Fator Azul permite manipular ou ampliar defesas baseadas em energia.
Essas habilidades normalmente interagem com o sistema de Escudos (Escudo = B × 3)
ou permitem novas formas de proteção energética.

Exemplos de termos:

- energia
- escudos
- absorção mágica
- defesas especiais

## Callback

Caso o mestre permita meta-habilidades ou magias e queira consequencias do uso decorrente.

Se ultrapassar 100%, o personagem desmaia.
O Mestre pode aplicar uma penalidade ou sequela permanente ao personagem segundo seus criterios.

```text
| Intensidade | Recuperação |
|-------------|-------------| leve | horas |
| moderado | dias |
| grave | semanas |
| extremo | meses |
```

## Absorção Especial

Absorção Especial permite reduzir dano recebido usando pontos de vetor.
Cada ponto utilizado do vetor reduz 1 ponto de dano recebido.
O mestre pode ajustar essa proporção dependendo do universo do jogo.

Absorve dano com pontos Vermelho (R).

Caso tenha fator Azul (B):

```text
Absorve dano com pontos Azul (B).
Absorve dano cortante, elemental e físico especial.
```

## Defesa Completa

Usuário se prepara para receber o golpe ao invés de se esquivar.
Anula tentativas de agarrar.

```text
Defesa Completa = **Defesa + Absorção Especial**
```

## Exemplos de Habilidades

Mobilidade

- voo
- levitação
- teleporte

Percepção

- visão noturna
- sensores avançados
- percepção ampliada

Tecnologia

- controle remoto de dispositivos
- habilidades de hacking
- interface com drones

## Como criar uma habilidade RGB

Toda habilidade deve definir:

- vetor utilizado
- custo (callback ou vetor)
- efeito
- duração
- limite por combate

Veja também:

- [Criação de Personagem](character_creation.md)
- [Atributos](attributes.md)

← [Voltar para README](README.md)
