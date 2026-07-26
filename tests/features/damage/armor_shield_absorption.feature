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

  Rule: Health consequences use exact thresholds

    Scenario: Damage equal to remaining health downs the target
      Given a target with current health 3, no armor and no shield
      When the target receives an impact of 3
      Then health must lose 3 points
      And the target must become downed
      And the target must not remain healthy

    Scenario: Fully mitigated damage does not injure
      Given a target with armor 5 and no shield
      When the target receives an impact of 4 with penetration 0
      Then health must lose 0 points
      And the target must remain healthy
