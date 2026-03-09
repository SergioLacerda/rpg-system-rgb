# Ficha de Personagem

A ficha de personagem resume os **valores principais de um personagem no Sistema RGB**.

Ela organiza os três vetores RGB junto com valores derivados e equipamentos.

## Informações Básicas

```text
Nome
Nível
```

## Vetores RGB

```text
R (Vermelho) → Poder / Dano
G (Verde) → Mobilidade / Reação
B (Azul) → Defesa / Energia
```

Esses vetores determinam a maioria das ações no sistema.

## Valores Derivados

Valores derivados são calculados a partir dos vetores RGB.

```text
Vida = 4 + (R × 2)
```

Escudo energético opcional (se o cenário permitir):

```text
Escudo Energético = B × 3
```

Sistema opcional de fadiga de habilidades:

```text
Callback (se utilizado no cenário da campanha)
```

## Equipamentos

```text
Equipamentos
Armas
Armadura
```

Os equipamentos interagem com os sistemas de combate e dano.

## Relação com o Sistema RGB

A ficha de personagem conecta diversos módulos do sistema RGB.

```text
Vetores → definem as capacidades do personagem
Vida → derivada de R
Escudo → derivado de B
Equipamentos → modificam a interação em combate
```

## Filosofia de Design

A ficha de personagem RGB segue três princípios:

- **estrutura mínima** – poucos números definem o personagem  
- **leitura rápida** – todos os valores principais visíveis rapidamente  
- **modularidade do sistema** – sistemas opcionais (escudo, callback) podem ser adicionados conforme a campanha  

Isso permite que a mesma ficha funcione em:

- campanhas modernas
- cenários de fantasia
- mundos de ficção científica
- jogos com superpoderes

## Exemplo Mínimo de Ficha

```text
Nome: Stormy
Nível: 3

R: 3
G: 3
B: 1

Vida: 10
Escudo: 3

Armadura: Armadura Média
Armas: Pistolas Duplas
Equipamentos: Equipamento Tático
```

← [Voltar para README](README.md)
