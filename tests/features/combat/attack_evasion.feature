Feature: Attack versus evasion
  RGB Core resolves direct pressure against evasion through deterministic margin.

  Rule: R applies pressure while G changes relation to that pressure

    Scenario Outline: Direct attack compared with evasion
      Given an attacker with R equal to <attack>
      And a defender with G equal to <evasion>
      When the attacker makes a direct attack
      Then the margin outcome must be <outcome>
      And damage must start only when the outcome is successful

      Examples:
        | attack | evasion | outcome                  |
        | 2      | 4       | failure_with_opportunity |
        | 3      | 3       | success_with_cost        |
        | 5      | 3       | success                  |
        | 7      | 3       | strong_success           |

  Rule: Margin bands follow the canonical table

    Scenario Outline: Canonical margin boundaries
      Given an attacker with R equal to <attack>
      And a defender with G equal to <evasion>
      When the attacker makes a direct attack
      Then the margin outcome must be <outcome>

      Examples:
        | attack | evasion | outcome                  |
        | 6      | 3       | strong_success           |
        | 5      | 3       | success                  |
        | 4      | 3       | success                  |
        | 3      | 3       | success_with_cost        |
        | 2      | 3       | failure_with_opportunity |
        | 1      | 3       | failure_with_opportunity |
        | 0      | 3       | clear_failure            |
        | 0      | 4       | clear_failure            |

  Rule: Modifiers shift the margin before classification

    Scenario: Situational modifier converts a miss into a hit
      Given an attacker with R equal to 2
      And a defender with G equal to 3
      And the attacker has a situational modifier of 2
      When the attacker makes a direct attack
      Then the margin must be 1
      And the margin outcome must be success

  Rule: The teaching shorthand stays consistent with margins

    Scenario: R greater than or equal to G always yields a successful outcome
      Given any attacker R from 0 to 10
      And any defender G from 0 to 10
      When R is greater than or equal to G
      Then the resolution outcome must be successful
