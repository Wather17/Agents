# Workflow de Trabalho

## Início da Sessão

Liste todos os arquivos com `TASK` no nome rodando `ls TASK-*.md` em ordem alfabética — essa é sua fila de execução.

Antes de qualquer coisa, leia o `TASK.md` — ele define o escopo completo do que você vai implementar. Se o `TASK.md` não existir ou estiver vago, anote em `BLOCKERS.md` e aguarde definição humana.

Com o escopo claro:

1. Rode `git status` — analise o estado do repositório. Se houver arquivos pendentes, analise as mudanças e commite antes de começar
2. Rode a bateria de testes completa — só avance se o repositório estiver estável. Se houver falhas pré-existentes, anote em `BLOCKERS.md` e passe para a próxima task
3. Crie uma branch específica para a task — uma task, uma branch

## Execução

Quebre a task em partes menores antes de escrever qualquer código. Pense sistematicamente: o que precisa existir antes do quê, quais são as dependências entre as partes.

Execute uma parte por vez, commitando atomicamente ao longo do caminho — não acumule tudo para o final.

### Testes

Toda implementação nova que pode ser testada deve ser testada. Ao escrever testes considere:

- **Caminho feliz** — o comportamento esperado no fluxo normal
- **Edge cases** — entradas vazias, nulas, limites, valores inesperados
- **Regressão** — garanta que o que funcionava antes continua funcionando
- **Integração** — se sua implementação interage com outros módulos, teste a fronteira

### Análise de Regressão

Antes de concluir, mapeie ativamente o que pode ter sido afetado pela sua mudança:

- Funções ou módulos que chamam o código que você alterou
- Tipos e interfaces que dependem das estruturas modificadas
- Comportamentos implícitos que o código anterior garantia
- Rode a bateria de testes completa e analise qualquer falha — não ignore warnings

## Critérios de Escalonamento

Durante a execução você vai encontrar dois tipos de problema:

**Você resolve sozinho:**
- Erros de sintaxe, linting, tipos
- Testes quebrando por causa da implementação atual
- Dependência faltando
- Refactor necessário para completar a task dentro do escopo definido

**Você anota em `BLOCKERS.md` e passa para a próxima task:**
- Decisão arquitetural não prevista na task
- Conflito com código existente que muda o escopo original
- Credencial, variável de ambiente ou acesso ausente
- Ambiguidade de requisito onde qualquer caminho tem trade-off significativo

Formato do blocker:

```markdown
## Blocker — [timestamp]
**O que você estava fazendo:** descrição da etapa em execução
**Problema encontrado:** descrição objetiva
**Decisão necessária:** o que precisa ser decidido ou fornecido
**Arquivos afetados:** lista dos arquivos relevantes
```

## Continuidade

Você nunca para ocioso. Se uma task encontrar um blocker de escopo humano:

1. Commite tudo que você fez até o momento na branch atual
2. Anote o blocker em `BLOCKERS.md` com o formato padrão
3. Passe para a próxima task do `TASK.md`

**Tasks no `TASK.md` são independentes entre si.** Se uma task depende do resultado de outra, elas devem estar agrupadas como uma task única ou como subtasks dentro da mesma entrada. Dependência entre tasks separadas é um erro de planejamento humano, não um problema para você resolver.

Você só para completamente quando:
- Todas as tasks foram concluídas
- Todas as tasks restantes estão bloqueadas aguardando decisão humana

## Conclusão da Task

Com a implementação completa e todos os testes passando:

- Certifique-se de que seus commits estão atômicos — mudanças relacionadas agrupadas, contexto independente separado
- Mensagens de commit descrevem o **porquê**, não só o **o quê**
- Rode a bateria de testes uma última vez antes do merge

## Merge

Faça o merge da branch da task para a `develop`. Após o merge:

```bash
git push origin develop
```

Limpe a branch da task se não for mais necessária.

## Ambiente
- Sistema: Windows
- Shell: PowerShell
- NUNCA usar && para encadear comandos — não funciona no PowerShell
- Para encadear: usar ; ou comandos separados
- Exemplo correto: `npm install; npm run build`
- Exemplo errado: `npm install && npm run build`