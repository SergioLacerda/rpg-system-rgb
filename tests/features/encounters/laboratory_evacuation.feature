Feature: Laboratory evacuation
  RGB Core encounters can be completed through objectives instead of defeating all opponents.

  Background:
    Given a researcher must reach the exit
    And mercenaries block the corridor
    And the reactor collapses after 4 rounds

  Scenario: The team wins by protecting and repositioning
    Given the warden protects the researcher
    And the runner opens an alternate route
    When the researcher reaches the exit before the fourth round
    Then the encounter objective must succeed
    And defeating every mercenary must not be required
    And no undefined rule step may be recorded

  Rule: Objectives can succeed or fail on a round deadline

    Scenario: Objective met before the deadline succeeds
      Given the laboratory evacuation encounter with a 4 round deadline
      When the researcher reaches the exit on round 3
      Then the encounter objective must succeed

    Scenario: Deadline reached without the objective fails the encounter
      Given the laboratory evacuation encounter with a 4 round deadline
      When round 4 ends without the researcher at the exit
      Then the encounter objective must fail
      And the failure reason must be the reactor collapse
