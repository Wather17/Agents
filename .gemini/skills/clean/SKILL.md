---
name: clean
description: Fecha a issue atual via GitHub CLI, marca o checklist e deleta arquivos .md temporários para evitar poluição do RAG.
---

Finalizamos a implementação e os testes. Agora precisamos fechar o ciclo para manter o repositório limpo.

**Regras de Execução:**
1. Verifique qual issue acabamos de resolver consultando o `PLANO_EXECUCAO.md`.
2. Use o GitHub CLI para fechar a issue correspondente. Comando: `gh issue close <numero_da_issue> -c "Resolvido via Antigravity."`
3. Atualize o arquivo `PLANO_EXECUCAO.md`, marcando a tarefa com um `[x]`.
4. Avalie o `PLANO_EXECUCAO.md`: Se TODAS as issues listadas estiverem concluídas, delete os arquivos temporários `PENDENCIAS.md` e `PLANO_EXECUCAO.md` usando os comandos de terminal apropriados.
5. Me confirme no chat qual issue foi fechada.
