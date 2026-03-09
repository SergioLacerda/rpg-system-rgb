# Modelo de Dano

O Sistema RGB utiliza um modelo de dano em camadas que combina
armas, armaduras, escudos e habilidades especiais.

Esta seção explica como o dano é aplicado durante o combate.

## Fluxo de Dano

A resolução do dano segue a seguinte sequência:

1. Verificação de ataque
2. Dano da arma
3. Interação com armadura
4. Interação com escudo
5. Dano restante aplicado ao personagem

## 1. Resolução de Ataque

Um ataque é bem-sucedido quando o atacante supera a capacidade
defensiva do defensor.

→ Veja [Ataque e Defesa](attack_and_defense.md)

## 2. Dano da Arma

Cada arma define um valor fixo de dano.

O dano da arma representa o impacto base antes das reduções defensivas.

As armas são descritas na seção de armas.

→ Veja [Armas](../weapons/README.md)

## 3. Proteção da Armadura

Armaduras fornecem proteção física.

A armadura reduz o dano recebido dependendo do tipo.

Categorias típicas de armadura:

Armadura Leve  
Armadura Média  
Armadura Pesada

→ Veja [Armaduras](../equipment/armor.md)

## 4. Proteção por Escudo

Alguns personagens podem possuir escudos.

Escudos absorvem dano antes que ele atinja o personagem.

Escudos de energia normalmente são calculados como:

```text
Escudo = B × 3
```

→ Veja [Escudos](../equipment/shields.md)

## 5. Penetração

Algumas armas podem ignorar parte da proteção da armadura.

Penetração representa a capacidade de uma arma atravessar armaduras.

→ Veja [Penetração](../weapons/penetration.md)

## 6. Dano Final

Após armadura e escudos serem resolvidos, o dano restante é aplicado
ao personagem.

Esse dano afeta a vitalidade do personagem.

## Habilidades Especiais

Algumas habilidades podem modificar a resolução do dano.

Exemplos:

- habilidades de absorção
- escudos de energia
- técnicas defensivas

→ Veja [Habilidades](../core/skills_and_abilities.md)

## Visão Geral da Interação de Dano

O modelo de dano RGB pode ser resumido como:

```text
Dano da Arma
→ Redução por Armadura
→ Absorção por Escudo
→ Dano Restante
```

Esse modelo em camadas permite que o sistema permaneça simples,
ao mesmo tempo que suporta decisões táticas de combate.
