# Exemplo Completo de Combate

Este exemplo demonstra um **combate completo** utilizando o Sistema RGB.

Dois personagens entram em combate dentro de um ambiente tático.

## Personagens

## Stormy

R = 3  
G = 3  
B = 2  

Valores derivados:

Vida = 4 + (R × 2) = 10  
Escudo = B × 3 = 6  
Movimento = G × 2 = 6 metros  

Equipamentos:

- Armadura Média
- Duas Pistolas

## Agente de Segurança

R = 2  
G = 2  
B = 1  

Valores derivados:

Vida = 8  
Escudo = 3  
Movimento = 4 metros  

Equipamentos:

- Armadura Leve
- SMG

## Situação Inicial

Os personagens começam a **6 metros de distância**.

Stormy possui cobertura parcial atrás de um veículo.

## Turno 1

## Movimento

Stormy se move 2 metros para melhorar a linha de visão.

Fórmula de movimento:

Movimento = G × 2

Stormy poderia se mover até **6 metros**, mas escolhe apenas um reposicionamento menor.

## Ataque

Stormy dispara uma das pistolas.

Dano da arma:

Dano = 4  
Penetração = 1

## Interação com Armadura

Armadura do alvo:

Armadura Leve = 2  
Penetração = 1  
Armadura Efetiva = 1

Dano após armadura:

4 − 1 = 3

## Interação com Escudo

Escudo do alvo:

Escudo = 3

O escudo absorve todo o dano.

Escudo Restante = 0  
Dano na Vida = 0

## Turno 2

O agente de segurança contra-ataca.

Arma:

Rajada de SMG  
Dano = 5  
Penetração = 1

Armadura de Stormy:

Armadura Média = 4  
Armadura Efetiva = 3

Cálculo do dano:

5 − 3 = 2

O escudo de Stormy absorve o dano.

Escudo Restante = 4  
Dano na Vida = 0

## Resumo do Combate

Este exemplo ilustra o modelo de dano do sistema RGB.

Fluxo de dano:

Arma  
↓  
Penetração  
↓  
Armadura  
↓  
Escudo  
↓  
Personagem

## Observação Tática

Stormy possui vantagem devido a:

- maior mobilidade (G)
- escudo mais forte (B)

O agente depende mais do dano da arma.

Isso demonstra como **diferentes distribuições de RGB influenciam a estratégia de combate**.

← [Voltar para README](README.md)
