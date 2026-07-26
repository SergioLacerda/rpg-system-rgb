Feature: State lifecycle ownership
  Every Core V2 state has at least one producing source and a defined way it
  is cleared, per decision matrix D-006 closure evidence.

  Rule: Every core state is reachable through a procedure

    Scenario: No orphan states exist in the core model
      Given the list of core states
      Then each state must declare at least one producing procedure
      And each state must declare how it is cleared

  Rule: Recovery restores continuity

    Scenario: Stabilizing an injured character
      Given an injured character with current health above zero
      When an ally sustains the character with B using the stabilize procedure
      Then the character must become stabilized
      And the character must no longer be injured
