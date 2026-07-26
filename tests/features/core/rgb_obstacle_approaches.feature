Feature: RGB approaches to obstacles
  RGB Core classifies obstacle handling by what changes in the situation.

  Rule: R transforms the pressure source

    Scenario: Break a blocked door
      Given a blocked door prevents passage
      When a character applies sufficient R pressure
      Then the obstacle must become vulnerable
      And the passage can become accessible through source change

  Rule: G changes relation to the pressure source

    Scenario: Bypass a blocked door
      Given a blocked door prevents passage
      And an alternate window route exists
      When a character repositions with G
      Then the character must become covered
      And the door must remain unchanged

  Rule: B preserves action under pressure

    Scenario: Hold a closing door
      Given a door is closing under pressure
      When a character sustains the passage with B
      Then the character must become guarded
      And the passage remains available temporarily
