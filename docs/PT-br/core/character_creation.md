# Criação de Personagem

A criação de personagens no **Sistema RGB** é intencionalmente simples.  
Os jogadores distribuem pontos entre os três vetores RGB para definir como seu personagem atua em combate e interação.

```text
        G
   mobilidade / posicionamento

R ---------------- B
poder / dano        escudo / defesa
```

- **R (Vermelho)** → poder ofensivo e força física  
- **G (Verde)** → mobilidade, reação e posicionamento  
- **B (Azul)** → resistência, energia e proteção defensiva  

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
Estilo                Foco de Vetor
--------------------  ----------------
Atacante pesado       Alto R
Combatente móvel      Alto G
Guardião defensivo    Alto B
```

Construções híbridas também são possíveis.

## Progressão do Personagem

À medida que os personagens ganham experiência, eles se tornam mais fortes.

```text
+2 pontos de vetor por nível
```

O Mestre determina quando os personagens sobem de nível de acordo com a campanha.

## Vida

A durabilidade do personagem é calculada usando o vetor Vermelho.

```text
Vida = 4 + (R × 2)
```

Exemplo:

```text
R = 3
Vida = 10
```

Valores maiores de **R** representam personagens mais fortes e resistentes.

## Referência Rápida

```text
Ataque     : R ≥ G
Movimento  : G × 2 metros
Vida       : 4 + (R × 2)
Escudo     : B × 3
```

Essas fórmulas resumem as principais interações mecânicas do sistema RGB.

## Filosofia de Design

A criação de personagens no RGB segue três princípios:

- **simplicidade** — poucos números definem o personagem  
- **diversidade tática** — diferentes distribuições de vetores criam estilos diferentes de jogo  
- **modularidade** — habilidades e equipamentos podem expandir o sistema sem alterar suas regras centrais  

Como todos os personagens distribuem pontos entre os mesmos três vetores, o sistema naturalmente produz construções equilibradas e variadas.

Veja também:

- [Atributos](attributes.md)
- [Combate](../combat/attack_and_defense.md)
- [Movimento](../combat/movement.md)

← [Voltar para README](README.md)
