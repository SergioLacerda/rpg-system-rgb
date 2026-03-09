# Visão Geral do Sistema

O Sistema RGB é construído em torno de três vetores fundamentais:

- **R (Vermelho)** — poder, ataques e dano
- **G (Verde)** — agilidade, mobilidade e reação
- **B (Azul)** — escudos, energia e resistência especial

Esses vetores definem como os personagens interagem com o mundo e com o combate.

## Estrutura do Sistema

O sistema é dividido em várias seções principais.

## Regras Básicas

As regras básicas definem os vetores RGB e a criação de personagens.

→ [Atributos](attributes.md)  
→ [Criação de Personagem](character_creation.md)  
→ [Habilidades](skills_and_abilities.md)

## Sistema de Combate

As regras de combate definem como os personagens interagem durante encontros.

→ [Ataque e Defesa](../combat/attack_and_defense.md)  
→ [Modelo de Decisão de Combate](../combat/combat_decision_model.md)  
→ [Modelo de Dano](../combat/damage_model.md)  
→ [Movimento](../combat/movement.md)

## Equipamentos

Equipamentos fornecem opções defensivas e táticas.

→ [Armaduras](../equipment/armor.md)  
→ [Escudos](../equipment/shields.md)  
→ [Equipamentos](../equipment/gear.md)

## Armas

Armas determinam como o dano é aplicado durante o combate.

→ [Armas de Fogo](../weapons/firearms.md)  
→ [Armas Corpo a Corpo](../weapons/melee.md)  
→ [Explosivos](../weapons/explosives.md)

## Referência

Exemplos e materiais de apoio para o jogo.

→ [Ficha de Personagem](../reference/character_sheet.md)  
→ [Exemplo de Combate](../reference/combat_example.md)  
→ [Exemplo de Jogo](../reference/gameplay_example.md)

## Fluxo de Jogo

Uma interação típica no jogo segue esta estrutura:

Criação de Personagem  
→ Seleção de Equipamentos  
→ Posicionamento Tático  
→ Resolução de Combate  
→ Recuperação e Progressão

## Notas Adicionais de Design

Explicações técnicas adicionais do sistema RGB podem ser encontradas aqui:

- [Arquitetura do Sistema RGB](../reference/rgb_system_architecture_notes.md)

## Modelo de Interação RGB

O sistema pode ser visualizado como um triângulo:

```text
        G
   mobilidade / reação

R ---------------- B
poder / dano    escudo / energia
```

Cada vetor representa uma estratégia diferente de interação no combate.
