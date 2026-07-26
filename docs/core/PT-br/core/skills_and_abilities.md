# Habilidades

Habilidades representam capacidades especiais além das ações normais do Sistema
RGB. Elas permitem que personagens expandam suas capacidades dependendo do
cenário da campanha.

Habilidades são **módulos opcionais** e podem representar:

- técnicas marciais
- poderes sobrenaturais
- tecnologia avançada
- sistemas de magia

O Mestre determina se essas habilidades estão disponíveis na campanha.

## Estrutura de Habilidades RGB

As habilidades no Sistema RGB são organizadas em torno dos **três vetores
centrais**.

```text
        G
   mobilidade / posicionamento

R ---------------- B
poder / dano      escudo / defesa
```

Cada vetor forma naturalmente uma **árvore de habilidades**.

- **R (Vermelho)** → habilidades de pressão, impacto, interrupção e
  transformação
- **G (Verde)** → habilidades de movimento, timing, evasão, alcance e
  posicionamento
- **B (Azul)** → habilidades de preservação, estabilização, bloqueio, escudo e
  resistência

Essa estrutura permite que o sistema escale naturalmente com diferentes
cenários.

## Árvore de Habilidades Vermelha (R)

Habilidades vermelhas ampliam **pressão ofensiva e força física**.

Habilidades típicas:

- Dobrar Potência
- Ataque Retardado
- Ataque de Pressão de Ar
- Golpe de Poder
- Quebra de Armadura

Essas habilidades aumentam a capacidade de um personagem causar impacto.

Exemplo de progressão:

```text
Nível 3 → +1 R em ataques à distância (Precisão)
Nível 6 → Golpe de poder duplo
Nível 9 → Ataque de impacto em área
```

## Árvore de Habilidades Verde (G)

Habilidades verdes focam em **mobilidade, velocidade e posicionamento tático**.

Habilidades típicas:

- Reflexos Ampliados
- Aceleração Momentânea
- Passo Fantasma (movimento curto similar a teleporte)
- Esquiva Automática
- Corrida Sobrenatural
- Salto Ampliado
- Movimento Silencioso
- Distorção de Movimento
- Evasão Sobrenatural

Essas habilidades permitem controlar espaço e evitar dano.

## Árvore de Habilidades Azul (B)

Habilidades azuis representam **manipulação de energia, endurance e sistemas
defensivos**.

Habilidades típicas:

- Escudo de Energia
- Escudo Atordoante
- Escudo Inteligente / Armadura Viva
- Invisibilidade
- Clones
- Absorção de Energia

Essas habilidades melhoram sobrevivência, continuidade e defesas baseadas em
energia.

## Contrato de Habilidade

Toda habilidade deve declarar os campos abaixo antes de poder ser tratada como
canônica, validada, empacotada em bundles ou usada por um Specialist:

| Campo | Propósito |
| --- | --- |
| `id` | identificador estável |
| `name` | nome usado em mesa |
| `vector` | vetor RGB primário |
| `tier` | nível, rank ou tier de acesso |
| `requirements` | pré-requisitos |
| `action_type` | ação, reação, passiva, postura ou timing especial |
| `cost` | pontos de vetor, escudo, callback, carga de item ou outro custo |
| `range` | alcance do alvo ou área |
| `duration` | instantânea, rodada, cena, persistente ou condicional |
| `effect` | efeito mecânico |
| `limits` | por turno, por combate, cooldown ou restrição narrativa |
| `tags` | tags de classificação |
| `source_status` | proposta, opcional, canônica, depreciada ou exemplo |

Habilidades sem esses campos são exemplos ou notas de design, não regras
canônicas.

## Fator Marcial (Módulo Opcional)

O **Fator Marcial** introduz técnicas defensivas avançadas.

Exemplos:

- Defesa Completa
- Absorção Especial
- técnicas de postura defensiva

O Mestre decide se essas técnicas exigem:

- uma ação
- uma reação
- uma postura defensiva

## Fator Azul (Módulo Opcional)

Em campanhas que permitem **magia ou meta-habilidades**, o vetor Azul pode
manipular defesas baseadas em energia.

Essas habilidades normalmente interagem com o sistema de escudo:

```text
Escudo = B × 3
```

Exemplos:

- escudos mágicos
- barreiras de energia
- absorção de energia
- campos defensivos avançados

## Sistema de Callback (Opcional)

Se a campanha inclui habilidades poderosas, o Mestre pode aplicar um **sistema
de callback** para representar custo ou fadiga causada por essas habilidades.

Se callback ultrapassar **100%**, o personagem desmaia.

```text
| Intensidade | Recuperação |
|-------------|-------------|
| leve        | horas |
| moderada    | dias |
| grave       | semanas |
| extrema     | meses |
```

O Mestre também pode aplicar consequências permanentes dependendo da situação.

## Absorção Especial

Absorção Especial permite que personagens reduzam dano recebido usando um
recurso declarado ou expressão vetorial.

```text
1 ponto de vetor gasto → reduz 1 dano
```

Usos possíveis:

- **Vermelho (R)** → absorção de impacto físico
- **Azul (B)** → absorção de energia ou dano especial

O Mestre pode ajustar essa proporção dependendo da campanha.

## Defesa Completa

Defesa Completa é um nome legado para bloqueio ou postura de guarda em que o
personagem se prepara para receber pressão em vez de esquivar.

Essa postura também anula tentativas de agarrar.

```text
Defesa Completa = procedimento de Bloqueio + efeito de absorção declarado
```

## Categorias de Habilidades de Exemplo

### Mobilidade

- voo
- levitação
- teleporte

### Percepção

- visão noturna
- sensores avançados
- percepção ampliada

### Tecnologia

- controle remoto de dispositivos
- habilidades de hacking
- interface com drones

## Criando Habilidades RGB

Toda habilidade deve usar o contrato acima. No mínimo, ela deve definir vetor,
custo, tipo de ação, efeito, duração, limites e status de origem.

Isso mantém habilidades consistentes com o Sistema RGB.

Veja também:

- [Criação de Personagem](character_creation.md)
- [Atributos](attributes.md)

← [Voltar para README](README.md)
