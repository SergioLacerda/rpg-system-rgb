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

  Rule: Starting characters distribute exactly 7 points

    Scenario: A valid starting distribution is accepted
      Given a starting character with R 3, G 2 and B 2
      Then the character must be valid

    Scenario: An oversized distribution is rejected for starting characters
      Given a starting character with R 5, G 2 and B 2
      Then character creation must fail with a budget error
