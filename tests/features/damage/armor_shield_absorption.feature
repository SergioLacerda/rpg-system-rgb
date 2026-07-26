Feature: Armor and shield absorption
  RGB Core applies damage in a fixed order before health consequences.

  Rule: Impact flows through penetration, armor, shield, and health consequence

    Scenario: Armor and shield reduce an incoming impact
      Given an incoming impact of 7
      And penetration of 1
      And target armor of 2
      And target shield reserve of 3
      When damage is resolved
      Then armor must reduce 1 point
      And shield must absorb 3 points
      And health must lose 3 points
      And the target must become injured

    Scenario: Absorption never creates negative damage
      Given an incoming impact of 1
      And target armor of 5
      And target shield reserve of 5
      When damage is resolved
      Then health must lose 0 points
      And shield reserve must not be negative
