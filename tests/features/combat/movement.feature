Feature: Movement distance
  Movement is governed by G on a 1 square = 1 meter grid.

  Rule: Movement distance is G times two meters

    Scenario Outline: Ground movement per turn
      Given a character with G equal to <g>
      When the character moves in one turn
      Then the maximum distance must be <meters> meters

      Examples:
        | g | meters |
        | 0 | 0      |
        | 3 | 6      |
        | 5 | 10     |

    Scenario: Aerial movement uses the same formula
      Given a flying character with G equal to 3
      When the character moves in one turn
      Then the maximum distance must be 6 meters
