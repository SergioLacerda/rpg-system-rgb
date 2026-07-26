Feature: Fixture parity
  YAML fixtures and Go fixtures describe the same characters, so authoring
  and test/simulation surfaces cannot silently drift apart.

  Rule: YAML fixtures and Go fixtures describe the same characters

    Scenario: Character vectors match across fixture sources
      Given the YAML fixture characters and the Go fixture characters
      Then both sets must contain the same character IDs
      And each character must have identical R, G and B values in both sources
