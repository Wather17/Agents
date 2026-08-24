# Agents

Workflow e tooling para agentes de IA em projetos de software.

Este repositório contém prompts, scripts e o CLI **`agent-init`**, usados para padronizar a forma como agentes de IA interagem com seus repositórios.

---

## agent-init

`agent-init` é um CLI escrito em Go que instala configurações de agentes de IA e o script de sincronização de issues do GitHub em qualquer repositório.

### Instalação

```bash
go install github.com/Wather17/Agents/agent-init@latest
```

O binário será instalado em `$GOPATH/bin`. Certifique-se de que esse diretório esteja no seu `PATH`.

### Comandos principais

| Comando | Descrição |
|---------|-----------|
| `agent-init` | Instala o template padrão (`GEMINI.md`) no repositório atual |
| `agent-init --agent opencode` | Instala o template do OpenCode (`AGENTS.md`) |
| `agent-init update` | Atualiza o CLI para a última versão |
| `agent-init upgrade` | Atualiza os prompts já instalados no repositório atual |
| `agent-init version` | Exibe a versão e os metadados do build |

Para mais detalhes, veja o [README completo do agent-init](agent-init/README.md).

---

## Templates suportados

| Template | Arquivo instalado | Convenção |
|----------|-------------------|-----------|
| `gemini` | `GEMINI.md` | Gemini CLI |
| `opencode` | `AGENTS.md` | OpenCode |

Os arquivos de prompt (`GEMINI.md`, `AGENTS.md`), o agente (`agents/issue-architect.md`) e as skills (`.agents/skills/refine-issues.md` e `.agents/skills/autonomous-batch.md`) são adicionados ao `.gitignore` do repositório alvo para manter o histórico limpo. O script `scripts/sync-issues.sh` é versionado.

---

## Estrutura do repositório

```
Agents/
├── .github/                    # CI e releases do agent-init
│   ├── ISSUE_TEMPLATE/         # Templates de issues autossuficientes
│   ├── PULL_REQUEST_TEMPLATE.md
│   └── workflows/
│       ├── ci.yml
│       └── release.yml
├── agent-init/              # Código-fonte do CLI
│   ├── README.md            # Documentação completa do agent-init
│   └── ...
├── agents/                  # Agentes especializados
│   └── issue-architect.md
├── .agents/skills/          # Skills de refinamento e execução
│   ├── refine-issues.md
│   └── autonomous-batch.md
├── scripts/                 # Scripts de automação
│   └── sync-issues.sh       # Sincroniza issues abertas do GitHub
├── GEMINI.md                # Prompt principal de persona e workflow
└── README.md                # Este arquivo
```

---

## Workflow

O workflow definido nos prompts segue um ciclo contínuo de lotes com PRs atômicos:

1. **Sincronizar** issues abertas com `scripts/sync-issues.sh`
2. **Escolher** uma issue da pasta `issues/`
3. **Criar** uma branch `issue/<numero>-<slug>`
4. **Desenvolver** e commitar com Conventional Commits
5. **Abrir** um Pull Request para `main`
6. **Fazer merge** e limpar o ambiente

---

## Mais informações

- Documentação completa do CLI: [`agent-init/README.md`](agent-init/README.md)
- Prompt principal: [`GEMINI.md`](GEMINI.md)
