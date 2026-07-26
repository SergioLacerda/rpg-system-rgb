Feature: Character creation derives canonical resources
  RGB Core V2 derives durability resources from R and B at character
  creation, per docs/core/en/core/character_creation.md.

  Rule: Derived formulas follow the canonical sheet

    Scenario Outline: Health and shield derivation
      Given a new character with R <r>, G <g> and B <b>
      Then max health must be <health>
      And max shield must be <shield>

      Examples:
        | r | g | b | health | shield |
        | 3 | 2 | 2 | 9      | 6      |
        | 2 | 2 | 3 | 9      | 9      |
        | 0 | 7 | 0 | 4      | 0      |

    Scenario: No preservation investment means no shield state
      Given a new character with R 4, G 3 and B 0
      Then the character must not have the shielded state
      And max shield must be 0
