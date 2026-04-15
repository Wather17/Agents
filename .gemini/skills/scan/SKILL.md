---
name: scan
description: Faz varredura de QA no código, identifica bugs/melhorias e cria GitHub Issues atômicas via CLI. Não altera código.
---

Atue como um Engenheiro de QA Sênior. Revise o código atual e identifique erros de sintaxe, bugs lógicos ou gargalos de otimização claros. 

**Regras de Execução:**
1. **NÃO escreva nem altere nenhum código agora.** O foco é apenas diagnóstico.
2. Liste os problemas encontrados em um "Plano de Implementação" temporário.
3. Para cada problema, execute o GitHub CLI (`gh issue create`) para registrar a pendência.
4. **Atomicidade:** Um problema = Uma issue. Nunca agrupe problemas diferentes.
5. **Estrutura da Issue:** - Título: Curto, técnico e direto ao ponto.
   - Corpo (`--body`): Descreva o que está acontecendo de errado, o motivo pelo qual é um problema e a sugestão técnica.
