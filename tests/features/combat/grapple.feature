Feature: Grapple
  A grapple succeeds only when the attacker matches mobility and pressure,
  per docs/core/en/combat/attack_and_defense.md.

  Rule: Grapple requires G >= G and R >= R

    Scenario Outline: Grapple resolution
      Given an attacker with G <ag> and R <ar>
      And a target with G <tg> and R <tr>
      When the attacker attempts a grapple
      Then the grapple must <result>

      Examples:
        | ag | ar | tg | tr | result  |
        | 3  | 3  | 3  | 3  | succeed |
        | 4  | 2  | 3  | 3  | fail    |
        | 2  | 4  | 3  | 3  | fail    |
        | 4  | 4  | 3  | 3  | succeed |
