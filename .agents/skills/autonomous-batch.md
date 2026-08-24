# Skill: Execução Autônoma em Lotes

Esta skill coordena a execução contínua de issues prontas. Um lote é uma fila de trabalho; cada issue continua isolada em sua própria branch, commits e PR.

## 1. Princípios

- `main` é a branch de integração.
- O GitHub é a fonte de verdade; `issues/` é apenas um cache local.
- Uma issue por branch e por PR é o padrão.
- O lote deve continuar até não haver mais issues prontas.
- Uma issue bloqueada não deve impedir as demais.
- Nenhuma mudança deve ser executada em paralelo no mesmo worktree.

## 2. Preparação do Lote

Antes de selecionar trabalho:

1. Execute `git status --short`.
2. Se houver alterações que não pertencem ao lote, preserve-as e pare.
3. Execute `git fetch --prune origin`.
4. Execute `git switch main` e `git pull --ff-only`.
5. Execute `./scripts/sync-issues.sh`.
6. Se a sincronização falhar, pare o lote e reporte o erro. Não trate falha de rede ou autenticação como fila vazia.
7. Leia todas as issues sincronizadas antes de escolher a primeira.

## 3. Classificação

Classifique cada issue:

- `status:ready`: contrato completo, sem bloqueios e com critérios verificáveis.
- `status:needs-refinement`: faltam contexto, decisão ou critério de aceite.
- `status:blocked`: depende de outra issue, acesso, decisão ou recurso indisponível.
- `status:in-progress`: já foi reivindicada por outro trabalho; não duplique.

Uma issue sem label de estado pode ser considerada pronta somente se passar o gate descrito em `refine-issues.md`.

Ordene issues prontas por:

1. Dependências satisfeitas.
2. `priority:critical`, `priority:high`, `priority:medium`, `priority:low`.
3. Menor risco e menor escopo.
4. Data de criação mais antiga.

Antes de iniciar, registre a claim com `status:in-progress` ou um comentário no GitHub quando os labels não estiverem disponíveis. Não reivindique uma issue que já tenha trabalho ativo sem verificar a branch e o PR correspondentes.

## 4. Execução de Cada Issue

Para a issue escolhida:

1. Leia corpo, comentários, dependências e critérios de aceite.
2. Verifique no repositório se os arquivos e símbolos descritos existem.
3. Consulte o histórico antes de decidir uma alteração relevante.
4. Crie `issue/<numero>-<slug>` a partir de `main`.
5. Implemente somente o escopo da issue.
6. Adicione testes para sucesso, erro e bordas relevantes.
7. Execute os comandos de validação da issue e os padrões do projeto.
8. Revise o diff e os secrets antes de commitar.
9. Faça commits pequenos e descritivos.
10. Abra o PR para `main` contendo `Resolves #<numero>`.

Não agrupe issues diferentes no mesmo PR. Se a implementação revelar trabalho adicional, registre uma nova issue e mantenha a atual limitada ao contrato aprovado.

## 5. Gates de Integração

O merge automático só é permitido quando todos os gates forem verdadeiros:

- Critérios de aceite atendidos.
- Testes, formatação, lint/vet e build locais aplicáveis passaram.
- Diff revisado e limitado ao escopo.
- PR sem conflitos.
- CI verde, verificada com `gh pr checks --watch`.
- Nenhuma revisão ou comentário pendente que exija mudança.

Com todos os gates atendidos, execute:

```bash
gh pr merge --merge --delete-branch
```

Depois do merge, atualize `main` com `git pull --ff-only`, remova a branch local integrada, confirme o fechamento da issue e execute `git fetch --prune origin`.

## 6. Bloqueios e Falhas

Se a issue estiver incompleta, registre a menor pergunta necessária e marque `status:needs-refinement` ou `status:blocked`. Continue para a próxima issue pronta.

Se um teste falhar por causa da implementação, corrija o código. Se a falha for preexistente ou externa ao escopo, confirme com histórico/testes, documente no PR e bloqueie apenas essa issue quando necessário.

Se houver conflito com alterações do usuário, não faça stash, reset ou checkout destrutivo. Preserve o worktree e peça orientação.

## 7. Encerramento do Lote

Continue até que não existam issues `status:ready`. Ao final, reporte de forma objetiva:

- Issues implementadas e PRs mergeados.
- Issues bloqueadas e o motivo exato.
- Issues que precisam de refinamento.
- Falhas de infraestrutura ou CI ainda pendentes.

Não declare o lote concluído enquanto houver uma issue pronta que possa ser executada com segurança.
