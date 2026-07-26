# Ataque e Defesa

## Iniciativa

Mais pontos em Verde (G) → age primeiro.
Ataque surpresa ignora iniciativa.

## Filosofia de Resolução

Core V2 resolve ataques com uma **margem**. Isso mantém o jogo determinístico
possível enquanto permite resultados mais ricos do que acerto ou erro plano.

```text
Margem de Ataque = Atacante (R + modificadores) - Defensor (G + modificadores)
```

Use a margem da seguinte forma:

| Margem | Resultado |
| ---: | --- |
| 3 ou mais | acerto forte |
| 1 a 2 | acerto |
| 0 | acerto com custo ou abertura para resposta |
| -1 a -2 | erro com oportunidade |
| -3 ou menos | erro claro |

O atalho de ensino antigo continua válido:

```text
R do Atacante ≥ G do Defensor
```

Esse atalho significa que o atacante obtém ao menos um acerto básico quando a
margem é zero ou maior.

## Ataque

Pressão ofensiva. O atacante usa Vermelho (R) para mudar o alvo ou a fonte de
perigo.

```text
margem = Atacante (R + modificadores) - Alvo (G + modificadores)
```

Se o resultado for um acerto, continue para o modelo de dano ou para a
consequência sem dano declarada.

## Procedimentos Defensivos

Defesa não é mais um número combinado. Escolha o procedimento que descreve o que
o defensor está fazendo.

| Procedimento | Dono Primário | Propósito |
| --- | --- | --- |
| Esquivar | G | evitar ou alterar contato |
| Reposicionar | G | mudar alcance, cobertura ou engajamento |
| Bloquear | B ou equipamento | receber pressão intencionalmente |
| Redução por armadura | equipamento | reduzir impacto por acerto |
| Absorção por escudo | recurso derivado de B ou equipamento | absorver dano restante |
| Proteger aliado | B, equipamento ou habilidade | redirecionar ou conter pressão |
| Interromper | R | parar uma ação aplicando pressão primeiro |
| Contra-atacar | R ou G por procedimento | responder pressão por força ou timing |

Não colapse esses procedimentos em uma soma genérica de armadura e escudo.
Armadura e escudo agem depois na pipeline de dano.

## Esquivar

```text
Esquivar usa G contra o R do atacante.
```

Esquivar muda contato, timing, alcance ou posição. Se a esquiva for bem-sucedida,
nenhum dano é aplicado, salvo se uma habilidade ou efeito de área disser o
contrário.

## Bloquear

Bloquear usa B, equipamento ou uma habilidade para receber pressão
intencionalmente. Um bloqueio pode reduzir a margem de ataque, proteger outro
alvo ou enviar o ataque para as camadas de armadura e escudo.

## Agarrar

Permite imobilizar o oponente ou movê-lo forçadamente durante um turno.

```text
A ação de agarrar tem êxito quando:

Atacante (G) ≥ Alvo (G)
Atacante (R) ≥ Alvo (R)

O atacante pode usar meta-habilidades.
O alvo pode usar bônus pré-declarados.
```

## Dano

Se um ataque causa dano, use o modelo de dano.

## Resolução de Dano

1. verificar acerto
2. estabelecer Fonte de Impacto
3. aplicar penetração
4. aplicar redução por armadura
5. aplicar absorção por escudo
6. aplicar dano restante à vida ou consequência de estado

## Meta Habilidades ou Magia

### Ataque mágico durante o combate

```text
A cada turno, pontos são recuperados para usar magias ou meta-habilidades.

Regeneração = 1 ponto B por nível
```

### Defesa mágica durante o combate

Defesa mágica é um escudo, bloqueio, resistência ou efeito de absorção tipado.
Ela deve declarar qual procedimento defensivo modifica e se age antes ou depois
da armadura.

Veja também:

- [Armas de Fogo](../weapons/categories/firearms.md)
- [Armaduras](../equipment/armor.md)

← [Voltar para README](../README.md)
