# Progressão

À medida que personagens ganham experiência, eles se tornam mais fortes.
A progressão no **Sistema RGB** é construída sobre um orçamento de
avanço fixo e pequeno por nível, que o jogador pode gastar como
crescimento bruto de vetor ou, quando a campanha usa o módulo de
Habilidades e Skills, como uma entre várias escolhas alternativas de
avanço.

Este documento é a fonte canônica para progressão de personagem. Ele
substitui o resumo que antes ficava embutido em
[Criação de Personagem](character_creation.md), que agora aponta para
cá.

## Orçamento de Avanço

```text
+2 pontos de vetor por nível
```

O Mestre determina quando os personagens sobem de nível, de acordo com a
campanha. Essa base nunca muda — cada escolha de avanço abaixo gasta o
mesmo orçamento fixo, não o soma.

## Escolhas de Avanço

Por padrão, o orçamento completo se torna pontos de vetor. Quando a
campanha permite o módulo [Habilidades e Skills](skills_and_abilities.md),
o jogador pode gastar o avanço de um nível em uma alternativa abaixo. Vale
exatamente uma escolha por nível; as escolhas não se acumulam dentro do
mesmo nível.

| Escolha | Efeito |
| --- | --- |
| **Crescimento de vetor** (padrão) | +2 pontos, distribuídos entre R, G e B como o jogador escolher |
| **Nova habilidade** | Ganha uma habilidade cujo `tier` e `requirements` (ver o Contrato de Habilidade em [Habilidades e Skills](skills_and_abilities.md)) o personagem atualmente cumpre |
| **Melhoria de habilidade** | Reduz o `cost` de uma habilidade conhecida, estende sua `duration`, ou aumenta o teto de `limits` em um degrau, a critério do Mestre |
| **Especialização** | Compromete-se com a árvore de habilidades de um vetor (Vermelho, Verde ou Azul): futuras escolhas de nova habilidade nessa árvore ignoram um tier de `requirements` |
| **Novo recurso** | Ganha um pool de recurso definido pela campanha (ex.: uma carga extra de reação, um buffer de recuperação de callback) dimensionado pelo Mestre |
| **Nova reação** | Ganha acesso a um procedimento de reação não disponível no tier atual do personagem |
| **Acesso a novo estado ou manobra** | Ganha a capacidade de declarar um estado tático ou manobra (ver [Ataque e Defesa](../combat/attack_and_defense.md), [Movimento](../combat/movement.md)) que antes era restrito ou indisponível |

## Especialização Na Prática

Especialização é a única escolha com efeito estrutural duradouro: ela
compromete um personagem com a identidade de um vetor em vez de conceder
um benefício imediato. Um personagem especializado em Vermelho não ganha
uma habilidade Vermelha só por se especializar — a especialização muda
como futuras habilidades Vermelhas são liberadas, conforme o campo
`requirements` do Contrato de Habilidade. Isso mantém a especialização
significativa sem exigir rastreamento de recursos além do que o Contrato
de Habilidade já define.

## Justificativa de Design

A análise inicial do RGB System V2 (intake `base_project`, §7.5) deixou
essa escolha em aberto: crescimento numérico puro, ou uma escolha por
avanço entre várias opções. Este documento resolve isso como um
**híbrido**, não uma invenção nova:

- A base numérica pura (`+2 pontos de vetor por nível`) já era canônica,
  declarada em Criação de Personagem antes deste documento existir. Este
  documento mantém isso inalterado como padrão — nenhum exemplo ou regra
  existente que assumia crescimento numérico puro é invalidado.
- [Habilidades e Skills](skills_and_abilities.md) já mostrava habilidades
  liberadas por nível em seus exemplos trabalhados (`Nível 3 → +1 R para
  ataques à distância`, `Nível 6 → Golpe de poder duplo`, `Nível 9 →
  Ataque de impacto em área`) — uma escolha de habilidade por avanço já
  estava implícita na prática, só nunca formalizada como opção
  documentada ao lado do crescimento de vetor. Este documento apenas
  torna esse padrão existente explícito e consistente com os campos
  `tier`/`requirements` do Contrato de Habilidade, que já existem
  especificamente para liberar acesso a habilidades por avanço.
- Um híbrido mantém intacta a filosofia de design declarada do sistema:
  criação de personagem e progressão favorecem **simplicidade** (um
  orçamento fixo, uma escolha por nível) e ao mesmo tempo suportam
  **diversidade tática** (a campanha, via o módulo opcional de
  Habilidades e Skills, pode oferecer mais que números puros sem
  adicionar novos subsistemas).

## Exemplo

Um personagem de nível 6 com `R=5, G=3, B=4` e o módulo de Habilidades e
Skills ativo pode escolher:

```text
Crescimento de vetor → R=7, G=3, B=4 (ou qualquer divisão de +2)
Nova habilidade      → ganha uma habilidade com tier <= tier atual do personagem
Especialização       → compromete-se com Vermelho, Verde ou Azul para futura liberação de habilidades
```

O jogador escolhe uma; a escolha não acumula nem carrega para o próximo
nível.

Veja também:

- [Criação de Personagem](character_creation.md)
- [Habilidades e Skills](skills_and_abilities.md)
- [Atributos](attributes.md)

← [Voltar para README](README.md)
