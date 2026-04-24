# Agent Workflow — Documentação

## Filosofia

O objetivo é simples: **deixar o agente trabalhar com autonomia de forma não perigosa**.

A ideia central é separar responsabilidades de forma clara:
- Você define o quê e o escopo
- O agente executa e toma decisões dentro do escopo
- Quando aparece algo além do escopo do agente, ele anota e segue
- Você volta, lê as notas, decide, e o ciclo recomeça
- O Agente deve trabalhar com autonomia mas que ele não execute comandos destrutivos para o projeto ou o ambiente de trabalho.

Isso evita dois problemas comuns: o agente interrompendo você pra perguntar coisas triviais, e o agente tomando decisões arquiteturais que não era pra ele tomar.

---

## Fluxo

```
Ideia → TASK.md → Agente executa → Erro de escopo humano? → BLOCKERS.md → Você analisa → Corrige → Repete
                                  ↓
                         Erro de escopo do agente → Resolve sozinho → Continua
```

## Para que serve cada arquivo

### GEMINI.md e CLAUDE.md
Lido automaticamente no início de cada sessão por ambos os CLIs. É a memória persistente do projeto entre sessões.

Coloque aqui:
- Stack e versões
- Convenções de código (nomenclatura, estrutura de pastas, padrões)
- O que o agente **nunca** deve fazer sem consultar você
- Como o workflow funciona (critérios de escalonamento, onde anotar blockers)
- Estado atual do projeto se relevante

Não coloque aqui:
- Detalhes da task atual (vai pro `TASK.md`)
- Erros e blockers (vão pro `BLOCKERS.md`)
- Contexto longo que muda toda sessão

---

### TASK.md
Escrito por você **antes** de abrir a sessão. A regra é simples: se você não escreveu o TASK.md, não abre a sessão.

Uma task bem definida elimina 90% das perguntas do agente. O modelo mínimo:

```markdown
## Objetivo
Descrição clara do que deve ser feito ao final.

## Escopo
- Arquivos e pastas relevantes
- Sistemas que podem ser tocados

## Não fazer
- O que está fora do escopo desta task
- Decisões que precisam de aprovação humana
```

Para tasks com dependência sequencial, use seções numeradas no mesmo arquivo.
Para tasks independentes, prefira sessões separadas para não acumular contexto desnecessário.

---

### BLOCKERS.md
Escrito pelo agente quando encontra um erro fora do seu escopo. Você lê quando voltar.

Formato padrão que o agente deve usar:

```markdown
## Blocker — [timestamp]
**O que estava fazendo:** descrição da task em execução
**Erro encontrado:** descrição objetiva do problema
**Decisão necessária:** o que você precisa decidir ou fornecer
**Arquivos afetados:** lista dos arquivos relevantes
```

---

## Configuração dos CLIs

### Claude Code

**Global** — `~/.claude/settings.json`
```json
{
  "defaultMode": "acceptEdits"
}
```

**Por repositório** — `.claude/settings.json`
```json
{
  "permissions": {
    "allow": [
      "Bash(git *)",
      "Bash(npm run *)",
      "Bash(npm install *)"
    ]
  },
  "hooks": {
    "PostToolUse": [
      {
        "matcher": "compact",
        "hooks": [{
          "type": "command",
          "command": "bash .agent/hooks/post-compact.sh"
        }]
      }
    ],
    "Notification": [
      {
        "matcher": "*",
        "hooks": [{
          "type": "command",
          "command": "bash .agent/hooks/notify.sh"
        }]
      }
    ]
  }
}
```

---

### Gemini CLI

**Global** — `~/.gemini/settings.json`
```json
{
  "defaultMode": "acceptEdits"
}
```

**Por repositório** — `.gemini/settings.json`
```json
{
  "hooks": {
    "AfterCompress": [
      {
        "hooks": [{
          "type": "command",
          "command": "bash .agent/hooks/post-compact.sh"
        }]
      }
    ],
    "Notification": [
      {
        "matcher": "*",
        "hooks": [{
          "type": "command",
          "command": "bash .agent/hooks/notify.sh"
        }]
      }
    ]
  }
}
```

---

## Critérios de Escalonamento

O que define se o agente resolve ou anota no BLOCKERS.md.

### Escopo do agente — resolve sozinho
- Erro de sintaxe, linting, tipos
- Teste quebrando por causa da implementação
- Dependência faltando
- Refactor necessário pra completar a task
- Qualquer erro que a solução esteja dentro do escopo definido na task

### Escopo humano — para e anota
- Decisão arquitetural não prevista na task
- Conflito com código existente que muda o escopo original
- Credencial, acesso ou variável de ambiente ausente
- Ambiguidade de requisito onde qualquer caminho tem trade-off significativo
- Qualquer coisa que exige contexto de negócio que o agente não tem

---

## Permissões

### Por que acceptEdits global
Editar arquivo é a operação mais comum e menos destrutiva com git presente. git checkout ou git reset desfaz qualquer edição. Não faz sentido confirmar escrita em todo repositório.

### Por que allow rules global
A minha stack de uso no dia a dia é bem definida, quando eu necessitar de algo novo eu devo conversar com o agente e modificar o rules global nos dois CLIs.

---

## Skill de Setup de Permissões

Crie em `.claude/skills/setup-permissions/SKILL.md`:

```markdown
---
name: setup-permissions
description: Analisa o projeto e configura .claude/settings.json com os allow rules adequados para o stack atual
---

Analise o projeto atual e gere um `.claude/settings.json` otimizado:

1. Leia o `package.json` (scripts, dependências)
2. Verifique se existe `Makefile`, `Dockerfile`, ou outros runners
3. Identifique o stack (framework, ORM, ferramentas de teste)
4. Gere os allow rules de Bash para os comandos mais prováveis
5. Escreva o arquivo `.claude/settings.json` mantendo o defaultMode herdado do global
```

Use com `/setup-permissions` ao iniciar um novo repositório.

---