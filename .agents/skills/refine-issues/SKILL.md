---
name: refine-issues
description: Refina ideias brutas, bugs e demandas técnicas em issues autossuficientes com critérios de aceite verificáveis. Use quando o usuário trouxer uma ideia vaga, ao auditar o código em busca de problemas ou ao preparar issues para execução autônoma.
---

# Skill: Refinamento de Issues Autossuficientes

Esta skill transforma ideias brutas, bugs encontrados e demandas técnicas em issues prontas para execução autônoma. Ela também pode ser usada para auditar o repositório e separar problemas em tarefas atômicas.

## 1. Quando Usar

Use esta skill quando:

- O usuário apresentar uma ideia, feature ou correção ainda vaga.
- Uma auditoria encontrar bugs, riscos, gargalos ou dívida técnica.
- Uma issue existente não tiver contexto suficiente para execução sem perguntas.

Não implemente código durante o refinamento, salvo se o usuário delegar explicitamente essa função. O resultado principal desta skill é uma especificação pronta e verificável.

## 2. Investigação Antes da Escrita

Antes de fazer perguntas ou propor uma solução:

1. Leia o código, testes e configuração relacionados.
2. Consulte o histórico com `git log`, `git blame` ou commits relevantes.
3. Procure padrões existentes e soluções semelhantes no repositório.
4. Identifique dependências, efeitos colaterais e limites de compatibilidade.
5. Separe fatos observados de hipóteses e decisões ainda abertas.

Não invente arquivos, símbolos, comandos ou causas. Quando o caminho exato não puder ser confirmado, registre isso como uma investigação necessária.

## 3. Perguntas de Prontidão

Faça perguntas curtas e desafiadoras somente sobre decisões que mudam a implementação:

- Qual comportamento atual reproduz o problema?
- Qual comportamento observável deve existir depois da mudança?
- Quais são os casos de erro, limites e concorrência?
- O que está explicitamente fora do escopo?
- Há impacto em compatibilidade, segurança, performance, dados ou custo externo?
- Quais dependências precisam estar prontas antes?
- Como cada critério será validado automaticamente ou manualmente?

Depois que as respostas forem obtidas, não mantenha questões abertas que alterem o comportamento. Questões residuais devem ser marcadas como bloqueio ou decisão não funcional.

## 4. Gate de Prontidão

Marque a issue como `status:ready` somente quando ela contiver contexto suficiente para outro agente executar sem nova entrevista. A issue está pronta quando possui:

- Contexto e problema.
- Estado atual, evidência e reprodução para bugs.
- Comportamento esperado.
- Critérios de aceite observáveis.
- Escopo e não-escopo.
- Solução proposta e restrições.
- Localização por arquivos e símbolos.
- Passos de implementação.
- Validação com comandos e resultados esperados.
- Dependências, riscos e compatibilidade.
- Nenhuma questão aberta essencial.
- Definition of Done.

Se o gate falhar, use `status:needs-refinement` ou `status:blocked`, registre a menor pergunta necessária e não entregue a issue ao executor.

## 5. Template Obrigatório de Issue

Use exatamente esta estrutura:

```markdown
# [Feature/Bug] Título claro e conciso

## 0. Metadados
- Tipo: Feature | Bug | Refactor | Docs
- Prioridade: critical | high | medium | low
- Escopo: [componente ou área]
- Tamanho: S | M | L
- Dependências: [issues ou "Nenhuma"]
- Bloqueia: [issues ou "Nenhuma"]
- Estado: status:ready | status:needs-refinement | status:blocked

## 1. Contexto e Problema
[Por que a mudança é necessária e qual problema técnico ou de negócio resolve.]

## 2. Comportamento Atual e Evidências
[Para bugs: reprodução, ambiente, logs e causa confirmada ou hipótese. Para features: comportamento existente e limitações.]

## 3. Comportamento Esperado
[Descrição observável, incluindo exemplos de entrada, saída e mensagens de erro quando aplicável.]

## 4. Critérios de Aceite
- [ ] Critério observável 1.
- [ ] Critério observável 2.
- [ ] Casos de erro e limites cobertos.

## 5. Escopo e Não-escopo
### Incluído
- [Mudanças que fazem parte desta issue.]

### Excluído
- [Mudanças que devem permanecer fora desta issue.]

## 6. Proposta de Solução e Restrições
[Abordagem recomendada, alternativas rejeitadas e restrições que não podem ser quebradas.]

## 7. Localização
- `caminho/arquivo.ext`: `SimboloOuFuncao`
- [Arquivos, testes, configurações e pontos de extensão afetados.]

## 8. Passo a Passo da Implementação
- [ ] Implementar [mudança concreta].
- [ ] Adicionar ou atualizar [teste concreto].
- [ ] Atualizar [documentação/configuração, se aplicável].

## 9. Validação
- Comando: `[comando exato]`
  - Resultado esperado: `[resultado]`
- Comando: `[comando de teste ou build]`
  - Resultado esperado: `[resultado]`

## 10. Dependências, Riscos e Compatibilidade
[Dependências, riscos de regressão, segurança, performance, migração e compatibilidade.]

## 11. Decisões e Questões Abertas
- Decisões tomadas: [lista]
- Questões abertas: Nenhuma questão que altere o comportamento.

## 12. Instrução de Autonomia
> Execute a issue de ponta a ponta após confirmar os arquivos no repositório. Busque contexto no código e no histórico antes de perguntar. Não expanda o escopo silenciosamente; crie uma nova issue para problemas descobertos fora dele.

## 13. Definition of Done
- [ ] Critérios de aceite atendidos.
- [ ] Testes de sucesso, erro e borda relevantes adicionados ou atualizados.
- [ ] Formatação, testes, lint/vet e build aplicáveis passaram.
- [ ] Documentação atualizada quando necessário.
- [ ] Diff revisado sem alterações fora do escopo ou secrets.
- [ ] PR aberto para `main`, CI verde e merge concluído.
```

## 6. Regras por Tipo

Para um bug, não aceite a issue sem comportamento atual, reprodução e evidência. Se a causa ainda for hipótese, transforme a confirmação da causa em passo explícito e não a apresente como fato.

Para uma feature, não aceite a issue sem casos de uso, comportamento esperado, critérios de aceite e não-escopo. Diferencie necessidade do usuário de uma implementação específica.

Para QA, crie uma issue por falha independente ou por grupo que compartilhe causa e correção. Informe severidade, impacto e evidência.

## 7. Criação no GitHub

Após consenso e gate de prontidão, prefira preservar a formatação com `--body-file`:

```bash
gh issue create \
  --title "[Feature/Bug] Título" \
  --body-file /caminho/para/issue.md
```

Se os labels existirem no repositório, aplique `priority:*`, `status:ready` e o label de área. O corpo da issue continua sendo a especificação; comentários servem para decisões e contexto posterior.
