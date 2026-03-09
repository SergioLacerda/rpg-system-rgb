# Sistema RGB — Regras em Uma Página

Uma referência mínima para começar a jogar o **Sistema de RPG RGB** imediatamente.

## Conceito Central

Os personagens são definidos por três vetores:

```text
| Vetor | Nome | Função |
|------|------|--------|
| R | Vermelho | dano, força, poder ofensivo |
| G | Verde | agilidade, mobilidade, reação |
| B | Azul | escudos, energia, defesa especial |
```

Os jogadores distribuem **7 pontos iniciais** entre R, G e B.

Exemplo:

R = 3  
G = 2  
B = 2

## Durabilidade do Personagem

```text
Vida = 4 + (R × 2)
Escudo = B × 3
```

Vida representa resistência física.  
Escudo representa energia ou proteção especial.

## Estrutura de Turno

A cada turno um personagem pode realizar:

```text
1 Ação
+
1 Ajuste Menor
```

Exemplos:

Ações:

- atacar
- mover
- defender
- usar habilidade

Ajustes menores:

- recarregar
- pequeno reposicionamento
- interagir com objeto

## Escolhas Táticas

A maioria das decisões de combate se divide em três categorias:

```text
Atacar → R
Mover → G
Defender → B
```

Isso cria três estratégias principais de combate:

```text
| Estratégia | Vetor |
|------------|------|
Atacante | R |
Escaramuçador | G |
Guardião | B |
```

Construções híbridas também são possíveis.

## Resolução de Ataque

Um ataque é bem-sucedido quando o atacante supera o defensor de acordo com as regras de combate.

Se o ataque for bem-sucedido, o dano é resolvido.

## Fluxo de Dano

```text
Teste de Ataque
↓
Dano da Arma
↓
Penetração
↓
Redução por Armadura
↓
Absorção por Escudo
↓
Dano Restante → Personagem
```

Definições:

Armadura → reduz dano por ataque  
Escudo → absorve dano acumulado  

## Exemplo de Dano

Dano da Arma: 7  
Penetração: 2  
Armadura: 4

```text
Armadura Efetiva = Armadura − Penetração
Armadura Efetiva = 2
Dano Final = 7 − 2 = 5
```

Se houver escudo, o dano é absorvido pelo escudo primeiro.

## Loop de Jogo

Uma sessão típica segue esta estrutura:

```text
Criação de Personagem
↓
Exploração
↓
Combate
↓
Recuperação
↓
Progressão da História
```

## Modelo Tático RGB

```text
        G
   mobilidade / reação

R - B
poder / dano    escudo / defesa
```

R → causar dano  
G → evitar dano  
B → absorver dano  

## Filosofia de Design

O sistema RGB foca em:

- relações numéricas simples
- decisões táticas de combate
- regras modulares
- cenários adaptáveis

As mesmas regras podem suportar:

- campanhas modernas
- mundos de fantasia
- ficção científica
- cenários com superpoderes
