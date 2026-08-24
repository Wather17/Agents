---
name: Bug report
about: Registrar um problema reproduzível com contexto suficiente para correção
title: "[Bug] "
labels: "status:needs-refinement"
---

# [Bug] Título claro e conciso

## 0. Metadados
- Tipo: Bug
- Prioridade: critical | high | medium | low
- Escopo: [componente ou área]
- Tamanho: S | M | L
- Dependências: [issues ou "Nenhuma"]
- Bloqueia: [issues ou "Nenhuma"]
- Estado: status:needs-refinement

## 1. Contexto e Problema

<!-- Explique por que o problema importa e qual impacto causa. -->

## 2. Comportamento Atual e Evidências

- Ambiente e versão:
- Pré-condições:
- Passos para reproduzir:
  1.
  2.
  3.
- Resultado atual:
- Logs, screenshots ou evidências:
- Causa: confirmada | hipótese | desconhecida

## 3. Comportamento Esperado

<!-- Descreva o resultado observável correto. -->

## 4. Critérios de Aceite
- [ ] O cenário reproduzido deixa de apresentar o problema.
- [ ] Casos de erro e limites relevantes estão cobertos.
- [ ] [Critério observável adicional]

## 5. Escopo e Não-escopo
### Incluído
-

### Excluído
-

## 6. Proposta de Solução e Restrições

<!-- Diferencie fatos, hipóteses e decisões. Registre alternativas rejeitadas. -->

## 7. Localização
- `caminho/arquivo.ext`: `SimboloOuFuncao`
-

## 8. Passo a Passo da Implementação
- [ ] Confirmar a causa ou executar a investigação descrita acima.
- [ ] Implementar a correção.
- [ ] Adicionar ou atualizar testes.
- [ ] Atualizar documentação, se aplicável.

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
