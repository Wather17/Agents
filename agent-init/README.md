# agent-init

CLI para instalar configuração de agentes de IA e o script de sincronização de issues em qualquer repositório.

## Instalação

```bash
go install github.com/Wather17/Agents/agent-init@latest
```

## Uso

Dentro do repositório que você quer configurar:

```bash
agent-init
```

Isso instala (com `--agent gemini`):
- `GEMINI.md` — prompt de persona e workflow (arquivo local, adicionado ao `.gitignore`)
- `scripts/sync-issues.sh` — script de sincronização de issues do GitHub
- Atualiza o `.gitignore` para ignorar `GEMINI.md` e `issues/`
- Cria um commit com Conventional Commit contendo `.gitignore` e `scripts/sync-issues.sh`

Com `--agent opencode`, instala `AGENTS.md` (convenção do OpenCode) no lugar de `GEMINI.md`.

## Flags

| Flag          | Padrão    | Descrição                                      |
|---------------|-----------|------------------------------------------------|
| `--agent`     | `gemini`  | Template do agente a ser instalado             |
| `--path`      | `.`       | Caminho do repositório alvo                    |
| `--force`     | `false`   | Sobrescreve arquivos existentes                |
| `--no-commit` | `false`   | Não cria o commit automaticamente              |

## Exemplos

```bash
# Instalar em outro diretório
agent-init --path /caminho/para/repo

# Usar o template do OpenCode
agent-init --agent opencode

# Sobrescrever arquivos existentes
agent-init --force

# Instalar sem criar commit
agent-init --no-commit
```
