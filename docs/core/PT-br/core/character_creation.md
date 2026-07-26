# Criação de Personagem

A criação de personagens no **Sistema RGB** é intencionalmente simples.
Jogadores distribuem pontos entre os três vetores RGB para definir como o
personagem atua em combate e interação.

```text
        G
   relação / posicionamento

R ---------------- B
pressão            preservação
```

- **R (Vermelho)** → pressão, força, ataque, disrupção e impacto
- **G (Verde)** → movimento, timing, evasão, alcance e posicionamento
- **B (Azul)** → endurance, bloqueio, estabilização, escudos e proteção

Essa distribuição define o estilo tático do personagem.

## Pontos Iniciais

Personagens começam com:

```text
7 pontos
```

Esses pontos devem ser distribuídos entre:

```text
R (Vermelho)
G (Verde)
B (Azul)
```

Exemplo:

```text
R = 3
G = 2
B = 2
```

Diferentes distribuições naturalmente criam diferentes papéis em combate.

```text
| Estilo | Foco de Vetor |
|--------|---------------|
Atacante pesado | R alto |
Combatente móvel | G alto |
Guardião defensivo | B alto |
```

Builds híbridas também são possíveis.

## Progressão do Personagem

À medida que personagens ganham experiência, eles se tornam mais fortes.

```text
+2 pontos de vetor por nível
```

O Mestre determina quando personagens sobem de nível de acordo com a campanha.

## Vida

A durabilidade do personagem usa tolerância à pressão e preservação.

```text
Vida = 4 + R + B
```

Exemplo:

```text
R = 3
B = 2
Vida = 9
```

Isso impede que **R** possua sozinho ataque e durabilidade. **R** maior ainda
ajuda personagens a suportar pressão física, enquanto **B** maior representa
preservação, estabilidade e endurance defensiva.

## Referência Rápida

```text
Ataque    : margem = R atacante - G defensor
Movimento : G × 2 metros
Vida      : 4 + R + B
Escudo    : B × 3
```

Essas fórmulas resumem as principais interações mecânicas do sistema RGB.

## Filosofia de Design

A criação de personagens no RGB segue três princípios:

- **simplicidade** — poucos números definem o personagem
- **diversidade tática** — diferentes distribuições de vetores criam estilos
  diferentes de jogo
- **modularidade** — habilidades e equipamentos podem expandir o sistema sem
  alterar suas regras centrais

Como todos os personagens distribuem pontos entre os mesmos três vetores, o
sistema naturalmente produz builds equilibradas e variadas.

Veja também:

- [Atributos](attributes.md)
- [Combate](../combat/attack_and_defense.md)
- [Movimento](../combat/movement.md)

← [Voltar para README](README.md)
