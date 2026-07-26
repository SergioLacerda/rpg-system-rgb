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
