Feature: Initiative ordering
  Initiative determines acting order within an encounter round, per
  docs/core/en/combat/attack_and_defense.md.

  Rule: Higher G acts first

    Scenario: The faster character opens the round
      Given a character "runner" with G equal to 5
      And a character "vanguard" with G equal to 2
      When the encounter round starts
      Then "runner" must act before "vanguard"

    Scenario: A surprise attack ignores initiative
      Given a character "vanguard" with G equal to 2 acting with surprise
      And a character "runner" with G equal to 5
      When the encounter round starts
      Then "vanguard" must act first
