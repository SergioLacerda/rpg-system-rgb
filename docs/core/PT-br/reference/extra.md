# Observação de Equilíbrio Matemático no Sistema RGB

Esta seção documenta uma propriedade de equilíbrio que surge ao cruzar as
principais fórmulas do sistema RGB. O objetivo é explicar como os vetores R, G e
B sustentam estratégias diferentes sem fazer um único vetor dominar ataque e
durabilidade.

## 1. Fórmulas Fundamentais do RGB

As regras básicas do sistema estabelecem as seguintes relações:

### Vida

```text
Vida = 4 + R + B
```

### Movimento

```text
Movimento = G × 2
```

### Escudo

```text
Escudo = B × 3
```

### Absorção Especial

Absorção física ou energética usa um recurso declarado pela habilidade ou
procedimento. Ela não substitui a pipeline de dano.

## 2. Impacto de Cada Vetor no Combate

Cada vetor influencia o combate de maneira diferente.

```text
| Vetor | Função |
|------|--------|
R | pressão, impacto e interrupção |
G | relação, mobilidade, evasão e posicionamento |
B | preservação, escudos, estabilização e proteção |
```

## 3. Escala de Resistência

Considerando um personagem com **5 pontos em cada vetor**:

```text
Vida = 4 + 5 + 5 = 14
Escudo = 5 × 3 = 15
```

Observação:

```text
Escudo ≈ Vida
```

Isso significa que a proteção baseada em energia possui resistência semelhante à
durabilidade física e de preservação, mas opera em outra camada.

## 4. Equilíbrio Entre os Vetores

Cada vetor contribui para sobrevivência de maneira diferente.

### Personagem focado em R

- maior pressão ofensiva
- mais tolerância a pressão física quando combinado à Vida

### Personagem focado em B

- Vida maior por preservação
- escudo energético mais forte
- melhor estabilização, bloqueio ou absorção quando declarados

Resultado:

Ambos podem alcançar resistência relevante, porém por mecânicas diferentes.

## 5. Papel do Vetor G

O vetor G não aumenta diretamente a resistência.

Em vez disso, ele melhora:

- evasão
- mobilidade
- timing
- posicionamento tático

Isso cria uma terceira forma de sobrevivência: evitar, alterar ou atrasar o
contato com a pressão.

## 6. Multiplicadores e Funções dos Vetores

```text
| Vetor | Expressão principal |
|------|----------------------|
R | margem de ataque e pressão |
G | movimento ×2 e relação com pressão |
B | Vida compartilhada e Escudo ×3 |
```

Isso cria três formas distintas de defesa:

```text
R → interromper ou encerrar a fonte de pressão
B → preservar continuidade e absorver pressão
G → evitar ou reposicionar a relação com a pressão
```

## 7. Arquétipos Naturais do Sistema

Sem criar classes formais, o sistema gera estilos naturais de personagem.

```text
| Estilo | Vetor dominante |
|------|------------------|
Atacante de pressão | R alto |
Evasivo | G alto |
Guardião de preservação | B alto |
```

## 8. Builds Híbridas

Exemplo:

R = 6  
B = 6

```text
Vida = 4 + 6 + 6 = 16
Escudo = 6 × 3 = 18
```

Resistência aproximada:

```text
16 + 18 = 34 pontos antes de consequências de estado
```

Isso demonstra que **combinações de vetores também são viáveis e competitivas**.

## 9. Triângulo de Equilíbrio do RGB

O sistema pode ser representado por um triângulo funcional:

```text
        G
   evasão / mobilidade

R ---------------- B
pressão            preservação
```

Cada vetor representa uma estratégia diferente para sobreviver ao combate.

## 10. Conclusão

O Sistema RGB apresenta um equilíbrio estrutural:

- R muda a fonte de pressão
- B preserva continuidade e sustenta escudos
- G reduz, altera ou evita contato com a pressão

Esse triângulo cria um sistema naturalmente balanceado sem exigir classes fixas
ou regras adicionais complexas.

← [Voltar para README](README.md)
