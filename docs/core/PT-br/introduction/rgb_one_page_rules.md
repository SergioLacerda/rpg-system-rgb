# Sistema RGB — Regras em Uma Página

Uma referência mínima para começar a jogar o **Sistema de RPG RGB** imediatamente.

## Conceito Central

Os personagens são definidos por três vetores:

```text
| Vetor | Nome | Função |
|------|------|--------|
| R | Vermelho | pressão, impacto, disrupção |
| G | Verde | relação, mobilidade, reação |
| B | Azul | preservação, escudos, estabilização |
```

Os jogadores distribuem **7 pontos iniciais** entre R, G e B.

Exemplo:

R = 3  
G = 2  
B = 2

## Durabilidade do Personagem

```text
Vida = 4 + R + B
Escudo = B × 3
```

Vida representa endurance física mais preservação.  
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
Pressionar   → R
Reposicionar → G
Sustentar    → B
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
Verificação de Acerto ou Contato
↓
Fonte de Impacto
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

Fonte de Impacto: 7  
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

R → causar dano ou mudar a fonte de pressão  
G → evitar dano ou mudar a relação com a pressão  
B → absorver dano ou preservar continuidade sob pressão  

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
