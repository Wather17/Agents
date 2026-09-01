# agent-init

CLI para instalar configuração de agentes de IA e o script de sincronização de issues do GitHub em qualquer repositório.

## Requisitos

- [Go](https://go.dev/doc/install) instalado (1.23+) para a instalação inicial via `go install`
- [GitHub CLI (gh)](https://cli.github.com/) instalado e autenticado para atualizar o CLI e usar o script `scripts/sync-issues.sh`

## Instalação

```bash
go install github.com/Wather17/Agents/agent-init@latest
```

O binário `agent-init` será instalado no `$GOPATH/bin`. Certifique-se de que esse diretório esteja no seu `PATH`.

## Atualizando

```bash
agent-init update
```

O comando baixa o binário da última GitHub Release correspondente ao sistema operacional e à arquitetura atual, valida o `checksums.txt` e substitui o executável instalado. Em Linux/WSL, a substituição é atômica.

Como alternativa manual ou bootstrap:

```bash
go install github.com/Wather17/Agents/agent-init@latest
```

## Uso

Dentro do repositório que você quer configurar, execute:

```bash
agent-init
```

Por padrão, o template `antigravity` é usado. `gemini` continua aceito como alias compatível. O comando instala:

- `GEMINI.md`, `AGENTS.md` ou `CLAUDE.md` — prompt de persona e workflow
- Skills com frontmatter no formato `<nome>/SKILL.md`
  - antigravity/gemini e codex: `.agents/skills/<nome>/SKILL.md`
  - opencode: `.opencode/skill/<nome>/SKILL.md` para as mesmas três skills (descoberta nativa pelo opencode)
  - claude: `.claude/skills/<nome>/SKILL.md`
- `scripts/sync-issues.sh` — script de sincronização de issues do GitHub
- `.agent-init.json` — manifesto local com schema versionado que registra os agentes instalados
- Atualiza o `.gitignore` para ignorar arquivos locais do agente e as issues sincronizadas
- Cria um commit com Conventional Commit

As três skills cobrem o ciclo completo: `refine-issues` (entrevista interativa), `audit-issues` (auditoria autônoma de QA) e `autonomous-batch` (execução da fila).

### O que é commitado

O commit gerado pelo CLI inclui:

- `.gitignore`
- `scripts/sync-issues.sh`

### O que fica ignorado

Os arquivos abaixo são instalados localmente, mas adicionados ao `.gitignore` para manter o repositório limpo:

- `GEMINI.md`, `AGENTS.md` ou `CLAUDE.md`
- Skills (`<skills>/<nome>/SKILL.md` conforme o template)
- `.agent-init.json`
- `issues/`

O `upgrade` remove automaticamente arquivos legados de versões anteriores, incluindo os flats `.agents/skills/*.md`, o agente `agents/issue-architect.md` e `.opencode/agent/issue-architect.md`.

### Templates disponíveis

| Template   | Arquivo instalado | Convenção    |
|------------|-------------------|--------------|
| `antigravity` | `GEMINI.md` | Antigravity |
| `gemini` | `GEMINI.md` | Alias de `antigravity` |
| `opencode` | `AGENTS.md` | OpenCode |
| `codex` | `AGENTS.md` | Codex |
| `claude` | `CLAUDE.md` | Claude Code |

Codex e OpenCode compartilham um `AGENTS.md` independente do host, mas mantêm skills em seus diretórios nativos. Para instalar vários agentes no mesmo repositório, execute o comando uma vez para cada alvo; o manifesto e os arquivos compartilhados são atualizados de forma idempotente.

## Flags

| Flag          | Padrão    | Descrição                                      |
|---------------|-----------|------------------------------------------------|
| `--agent`     | `antigravity` | Template do agente a ser instalado          |
| `--path`      | `.`       | Caminho do repositório alvo                    |
| `--force`     | `false`   | Sobrescreve arquivos existentes                |
| `--no-commit` | `false`   | Não cria o commit automaticamente              |

## Comandos

| Comando             | Descrição                                    |
|---------------------|-----------------------------------------------|
| `agent-init`        | Instala o template padrão no repo atual       |
| `agent-init update` | Atualiza o CLI para a última versão           |
| `agent-init upgrade`| Atualiza os prompts já instalados no repo atual |
| `agent-init version`| Exibe a versão e os metadados do build        |

### `agent-init upgrade`

Use esse comando quando os templates de prompt (`GEMINI.md`, `AGENTS.md`) forem atualizados no repo `Agents` e você quiser propagar essas mudanças para seus projetos:

```bash
agent-init upgrade
```

O que ele faz:
- Tenta atualizar o CLI pela última GitHub Release e reexecuta o comando usando o binário atualizado
- Atualiza os prompts e arquivos auxiliares já existentes no repo
- Instala as skills ausentes dos agentes registrados em `.agent-init.json` quando o prompt correspondente existe
- Não instala um novo template de agente (`GEMINI.md`, `AGENTS.md` ou `CLAUDE.md`)
- Não sobrescreve `scripts/sync-issues.sh` (para isso, use `agent-init --force`)
- Atualiza o `.gitignore` e cria um commit se for um repo git

Na primeira atualização de uma instalação antiga sem manifesto, `GEMINI.md` é migrado como `antigravity` e `AGENTS.md` como `opencode`, preservando a interpretação das versões anteriores.

Se a atualização automática do CLI falhar, o comando exibe um aviso e continua usando os templates embutidos na versão em execução.

## Releases

As releases do `agent-init` usam Semantic Versioning e tags no formato `agent-init/vMAJOR.MINOR.PATCH`, por exemplo:

```bash
git tag -a agent-init/v0.1.0 -m "release: agent-init v0.1.0"
git push origin agent-init/v0.1.0
```

O workflow em `.github/workflows/release.yml` compila os binários Linux AMD64, Linux ARM64 e Windows AMD64, publica uma GitHub Release e gera o `checksums.txt` usado pelo comando `update`.

## Workflow Autônomo

O arquivo de prompt e a skill `autonomous-batch` orientam a IA a:

- Ler e classificar todas as issues após uma sincronização bem-sucedida.
- Priorizar issues prontas por dependências, prioridade, risco e escopo.
- Executar uma issue por vez, mantendo branches, commits e PRs atômicos.
- Fazer merge em `main` somente depois dos critérios de aceite, validação local e CI verde.
- Bloquear issues incompletas sem interromper as demais do lote.
- Continuar até não haver mais issues prontas.

## Exemplos

```bash
# Instalar o template padrão (antigravity) no diretório atual
agent-init

# Instalar o template do OpenCode
agent-init --agent opencode

# Instalar Codex e Claude Code no mesmo repositório
agent-init --agent codex
agent-init --agent claude

# Instalar em outro diretório
agent-init --path /caminho/para/repo

# Sobrescrever arquivos existentes
agent-init --force

# Instalar sem criar commit automaticamente
agent-init --no-commit

# Atualizar o CLI
agent-init update

# Exibir a versão instalada
agent-init version

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
