---
name: Feature request
about: Descrever uma funcionalidade pronta para refinamento e execução autônoma
title: "[Feature] "
labels: "status:needs-refinement"
---

# [Feature] Título claro e conciso

## 0. Metadados
- Tipo: Feature
- Prioridade: critical | high | medium | low
- Escopo: [componente ou área]
- Tamanho: S | M | L
- Dependências: [issues ou "Nenhuma"]
- Bloqueia: [issues ou "Nenhuma"]
- Estado: status:needs-refinement

## 1. Contexto e Problema

<!-- Explique a necessidade, o usuário afetado e o problema que será resolvido. -->

## 2. Comportamento Atual e Evidências

<!-- Descreva o fluxo atual, limitações e evidências que justificam a mudança. -->

## 3. Comportamento Esperado

<!-- Descreva o fluxo esperado com exemplos de entrada, saída e erros. -->

## 4. Critérios de Aceite
- [ ] [Comportamento observável principal]
- [ ] [Caso de uso ou cenário alternativo]
- [ ] Casos de erro e limites relevantes estão cobertos.

## 5. Escopo e Não-escopo
### Incluído
-

### Excluído
-

## 6. Proposta de Solução e Restrições

<!-- Descreva a abordagem recomendada, alternativas rejeitadas e restrições. -->

## 7. Localização
- `caminho/arquivo.ext`: `SimboloOuFuncao`
-

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

<!-- Informe impactos de regressão, segurança, performance, dados e compatibilidade. -->

## 11. Decisões e Questões Abertas
- Decisões tomadas:
- Questões abertas:

## 12. Instrução de Autonomia

> Execute a issue de ponta a ponta após confirmar os arquivos no repositório. Busque contexto no código e no histórico antes de perguntar. Não expanda o escopo silenciosamente.

## 13. Definition of Done
- [ ] Critérios de aceite atendidos.
- [ ] Testes de sucesso, erro e borda relevantes passaram.
- [ ] Formatação, testes, lint/vet e build aplicáveis passaram.
- [ ] Documentação atualizada quando necessário.
- [ ] Diff revisado sem alterações fora do escopo ou secrets.
- [ ] PR aberto para `main`, CI verde e merge concluído.
