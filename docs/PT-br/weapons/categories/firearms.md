# Armas de Fogo

Armas de fogo representam as **principais armas de combate à distância** em cenários modernos e táticos do sistema RGB.

Diferente de armas corpo a corpo, armas de fogo **não adicionam o vetor R (Vermelho) do usuário ao dano**.  
O dano é determinado pela própria arma, representando o poder balístico e a munição utilizada.

Armas de fogo interagem com o **Sistema de Combate RGB** e seguem o **Modelo de Dano** padrão do sistema.

## Estatísticas da Arma

Cada arma de fogo é definida por quatro atributos principais:

```text
Dano        → dano base causado pela arma
Capacidade  → número de munições no carregador
Recarga     → ações necessárias para recarregar
Tipo        → categoria da arma
```

## Tabela de Armas de Fogo

```text
Tipo                Dano   Capacidade   Recarga
------------------  -----  ----------   -------------Pistola Compacta    4      12           ação parcial
Pistola             4      15           ação parcial
Revólver            5      6            ação principal
SMG                 5      30           ação principal
Rifle de Assalto    7      30           ação principal
Rifle Pesado / DMR  8      20           ação principal
Rifle de Precisão   9      10           ação principal
Metralhadora        7      100          2 ações
```

## Papéis das Armas

Diferentes armas cumprem papéis táticos distintos.

### Pistolas

- armas leves e secundárias
- recarga rápida
- eficientes em curta distância

### SMGs

- alta cadência de tiro
- eficazes em curta e média distância

### Rifles de Assalto

- equilíbrio entre dano e capacidade
- eficazes na maioria das situações de combate

### Rifles Pesados / Precisão

- dano elevado
- menor capacidade de munição
- utilizados para funções especializadas

### Metralhadoras

- fogo de supressão contínuo
- grande capacidade de munição
- tempo de recarga maior

## Interação com o Sistema RGB

Armas de fogo interagem indiretamente com os vetores RGB.

```text
R → melhora capacidade ofensiva ou controle da arma
G → influencia posicionamento e mobilidade em combate
B → fornece proteção defensiva através de escudos
```

Entretanto, o **dano da arma vem da própria arma**, e não do valor de R do personagem.

## Interação com o Modelo de Dano

Quando um ataque com arma de fogo atinge um alvo, o fluxo de dano do sistema RGB é aplicado:

```text
Dano da Arma
↓
Penetração (se aplicável)
↓
Redução por Armadura
↓
Absorção por Escudo
↓
Dano Restante → Personagem
```

Isso mantém o combate consistente com o restante do sistema.

## Filosofia de Design

Armas de fogo no RGB seguem três princípios:

- **clareza** — estatísticas simples de entender
- **variedade tática** — armas diferentes cumprem funções diferentes
- **compatibilidade com o sistema** — integração direta com o modelo de dano RGB

Isso permite que armas de fogo sejam poderosas sem tornar o sistema complexo.

## Veja Também

- [Combate](../../combat/attack_and_defense.md)
- [Alcance de Armas](../mechanics/weapon_ranges.md)
- [Modelo de Dano](../../combat/damage_model.md)

← [Voltar para README](../README.md)
