# agent-init

CLI para instalar configuração de agentes de IA e o script de sincronização de issues do GitHub em qualquer repositório.

## Requisitos

- [Go](https://go.dev/doc/install) instalado (1.23+)
- Para usar o script `scripts/sync-issues.sh`, é necessário o [GitHub CLI (gh)](https://cli.github.com/) autenticado

## Instalação

```bash
go install github.com/Wather17/Agents/agent-init@latest
```

O binário `agent-init` será instalado no `$GOPATH/bin`. Certifique-se de que esse diretório esteja no seu `PATH`.

## Atualizando

```bash
agent-init update
```

Ou, manualmente:

```bash
go install github.com/Wather17/Agents/agent-init@latest
```

## Uso

Dentro do repositório que você quer configurar, execute:

```bash
agent-init
```

Por padrão, o template `gemini` é usado. O comando instala:

- `GEMINI.md` — prompt de persona e workflow
- `agents/issue-architect.md` — agente especializado (Issue Architect & Quality Guard)
- `.agents/skills/refine-issues.md` — skill de refinamento de ideias e criação de issues
- `scripts/sync-issues.sh` — script de sincronização de issues do GitHub
- Atualiza o `.gitignore` para ignorar arquivos locais do agente e as issues sincronizadas
- Cria um commit com Conventional Commit

### O que é commitado

O commit gerado pelo CLI inclui:

- `.gitignore`
- `scripts/sync-issues.sh`

### O que fica ignorado

Os arquivos abaixo são instalados localmente, mas adicionados ao `.gitignore` para manter o repositório limpo:

- `GEMINI.md` (ou `AGENTS.md` para o template `opencode`)
- `agents/issue-architect.md`
- `.agents/skills/refine-issues.md`
- `issues/`

### Templates disponíveis

| Template   | Arquivo instalado | Convenção    |
|------------|-------------------|--------------|
| `gemini`   | `GEMINI.md`       | Gemini CLI   |
| `opencode` | `AGENTS.md`       | OpenCode     |

## Flags

| Flag          | Padrão    | Descrição                                      |
|---------------|-----------|------------------------------------------------|
| `--agent`     | `gemini`  | Template do agente a ser instalado             |
| `--path`      | `.`       | Caminho do repositório alvo                    |
| `--force`     | `false`   | Sobrescreve arquivos existentes                |
| `--no-commit` | `false`   | Não cria o commit automaticamente              |

## Comandos

| Comando             | Descrição                                    |
|---------------------|-----------------------------------------------|
| `agent-init`        | Instala o template padrão no repo atual       |
| `agent-init update` | Atualiza o CLI para a última versão           |
| `agent-init upgrade`| Atualiza os prompts já instalados no repo atual |

### `agent-init upgrade`

Use esse comando quando os templates de prompt (`GEMINI.md`, `AGENTS.md`) forem atualizados no repo `Agents` e você quiser propagar essas mudanças para seus projetos:

```bash
agent-init upgrade
```

O que ele faz:
- Atualiza o CLI (`agent-init update`)
- Atualiza os prompts e arquivos auxiliares já existentes no repo
- Instala os arquivos auxiliares ausentes (`agents/issue-architect.md` e `.agents/skills/refine-issues.md`) quando `GEMINI.md` ou `AGENTS.md` já existe
- Não instala um novo template de agente (`GEMINI.md` ou `AGENTS.md`)
- Não sobrescreve `scripts/sync-issues.sh` (para isso, use `agent-init --force`)
- Atualiza o `.gitignore` e cria um commit se for um repo git

Se a atualização automática do CLI falhar, o comando exibe um aviso e continua usando os templates embutidos na versão em execução.

## Exemplos

```bash
# Instalar o template padrão (gemini) no diretório atual
agent-init

# Instalar o template do OpenCode
agent-init --agent opencode

# Instalar em outro diretório
agent-init --path /caminho/para/repo

# Sobrescrever arquivos existentes
agent-init --force

# Instalar sem criar commit automaticamente
agent-init --no-commit

# Atualizar o CLI
agent-init update

# Atualizar os prompts instalados no repo atual
agent-init upgrade
```

## Repositórios que não são git

Se o diretório alvo não for um repositório git, o CLI ainda instala os arquivos e atualiza o `.gitignore`, mas pula a etapa de commit. Use `--no-commit` para forçar esse comportamento mesmo em repositórios git.

## Adicionando novos templates

Para adicionar um novo template de agente:

1. Adicione o arquivo de prompt em `internal/templates/files/` (ex: `CLAUDE.md`).
2. Declare o novo agente em `internal/templates/templates.go`.
3. Implemente os retornos de `FilesFor()` e `IgnoredEntries()` para o novo agente.
4. Adicione testes em `internal/templates/templates_test.go` e `internal/installer/installer_test.go`.
5. Atualize este README com o novo template.
