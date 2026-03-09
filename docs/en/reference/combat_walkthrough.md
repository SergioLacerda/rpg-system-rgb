# Combat Walkthrough

This example demonstrates a **complete combat encounter** using the RGB System.

Two characters engage in combat inside a tactical environment.

## Characters

## Stormy

R = 3  
G = 3  
B = 2  

Derived values:

Health = 4 + (R × 2) = 10  
Shield = B × 3 = 6  
Movement = G × 2 = 6 meters  

Equipment:

- Medium Armor
- Dual Pistols

## Security Agent

R = 2  
G = 2  
B = 1  

Derived values:

Health = 8  
Shield = 3  
Movement = 4 meters  

Equipment:

- Light Armor
- SMG

## Initial Situation

The characters begin **6 meters apart**.

Stormy has partial cover behind a vehicle.

## Turn 1

## Movement

Stormy moves 2 meters to improve line of sight.

Movement formula:

Movement = G × 2

Stormy could move up to **6 meters**, but chooses a smaller reposition.

## Attack

Stormy fires one pistol.

Weapon damage:

Damage = 4  
Penetration = 1

## Armor Interaction

Target armor:

Light Armor = 2  
Penetration = 1  
Effective Armor = 1

Damage after armor:

4 − 1 = 3

## Shield Interaction

Target shield:

Shield = 3

Shield absorbs all damage.

Remaining Shield = 0  
Health Damage = 0

## Turn 2

The security agent retaliates.

Weapon:

SMG Burst  
Damage = 5  
Penetration = 1

Stormy armor:

Medium Armor = 4  
Effective Armor = 3

Damage:

5 − 3 = 2

Stormy shield absorbs damage.

Remaining Shield = 4  
Health Damage = 0

## Combat Summary

This walkthrough illustrates the RGB combat model.

Damage flow:

Weapon  
↓  
Penetration  
↓  
Armor  
↓  
Shield  
↓  
Character

## Tactical Observation

Stormy benefits from:

- higher mobility (G)
- stronger shield (B)

The agent relies more on weapon damage.

This shows how **different RGB builds influence combat strategy**.

← [Back to README](README.md)
