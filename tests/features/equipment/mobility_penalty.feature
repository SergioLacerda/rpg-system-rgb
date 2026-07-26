Feature: Equipment mobility penalty
  Armor and physical shields reduce effective G, per
  docs/core/en/equipment/armor.md and docs/core/en/equipment/shields.md.

  Rule: Each armor category has a fixed mobility penalty

    Scenario Outline: Armor reduces effective G
      Given a character with G equal to 6
      And <armor> armor equipped
      And no shield equipped
      When effective movement G is computed
      Then effective G must be <effective_g>

      Examples:
        | armor  | effective_g |
        | none   | 6           |
        | light  | 5           |
        | medium | 4           |
        | heavy  | 3           |

  Rule: Each physical shield category has a fixed mobility penalty

    Scenario Outline: Shield reduces effective G
      Given a character with G equal to 6
      And no armor equipped
      And <shield> shield equipped
      When effective movement G is computed
      Then effective G must be <effective_g>

      Examples:
        | shield | effective_g |
        | none   | 6           |
        | light  | 6           |
        | medium | 5           |
        | heavy  | 4           |

  Rule: Combined penalties never drive effective G below zero

    Scenario: Heavy armor and heavy shield together floor at zero
      Given a character with G equal to 2
      And heavy armor equipped
      And heavy shield equipped
      When effective movement G is computed
      Then effective G must be 0
