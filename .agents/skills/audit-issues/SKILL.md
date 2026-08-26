---
name: audit-issues
description: Audita o repositório de forma autônoma em busca de bugs, riscos de segurança e dívida técnica, criando issues atômicas autossuficientes sem diálogo. Use quando o usuário pedir revisão de qualidade proativa, auditoria do código ou varredura por problemas.
---

# Skill: Auditoria Autônoma do Repositório

Esta skill executa uma auditoria de qualidade no repositório e transforma os achados em GitHub Issues atômicas e prontas para execução autônoma. Ela não conduz entrevista com o usuário e não corrige código.

## 1. Quando Usar

Use esta skill quando:

- O usuário pedir revisão de qualidade proativa, auditoria ou varredura por bugs.
- Antes de releases importantes, como verificação adicional à CI.
- O usuário relatar sintomas vagos ("está lento", "está instável") que exijam investigação ampla.

Não use para ideias novas ou features: nesse caso, a entrevista da skill `refine-issues` é o caminho adequado.

## 2. Processo de Auditoria

1. Mapeie a estrutura: módulos, pontos de entrada, configuração, CI e scripts.
2. Priorize áreas de maior risco: código que manipula dados sensíveis, concorrência, I/O, autenticação e parsing de entrada externa.
3. Leia o histórico recente com `git log` em busca de hotfixes recorrentes, que indicam áreas frágeis.
4. Analise os testes: cobertura real dos fluxos críticos, testes que nunca falham e asserções vazias.
5. Inspecione padrões conhecidos de defeito: erros ignorados, recursos sem fechamento, condições de corrida, validação ausente, secrets expostos e dependências com vulnerabilidade conhecida.
6. Para cada suspeita, confirme a evidência antes de registrar: leia o código até entender o mecanismo do problema, não apenas o sintoma.

Diferencie rigorosamente:

- **Fato**: demonstrável no código ou reproduzível.
- **Hipótese**: plausível, mas não confirmada; registre como investigação necessária dentro da issue, nunca como causa afirmada.
- **Preferência estilística**: não vira issue; anote no relatório final apenas se for relevante.

## 3. Critérios para Abrir Issue

Abra uma issue somente quando o achado tiver:

- Mecanismo compreendido, não apenas sintoma.
- Impacto concreto: corrupção de dados, falha funcional, segurança, performance mensurável ou dívida que bloqueia evolução.
- Localização precisa por arquivos e símbolos.
- Correção razoavelmente delimitável: uma issue por problema, ou grupo de problemas que compartilhem a mesma causa.

Não abra issues para estilo, nomenclatura, opiniões arquiteturais sem impacto demonstrável ou problemas já registrados no repositório.

## 4. Formato das Issues

Cada issue segue obrigatoriamente o template estrito e o gate de prontidão definidos na skill `refine-issues`:

- gemini: `.agents/skills/refine-issues/SKILL.md`
- opencode: `.opencode/skill/refine-issues/SKILL.md`

Adaptações para o contexto de auditoria:

- A seção de comportamento atual deve conter evidência técnica: trecho, caminho, símbolo e cenário de disparo.
- Quando a causa ainda for hipótese, transforme a confirmação em passo explícito de implementação e registre-a como hipótese na seção de decisões, mantendo os critérios de aceite verificáveis.

## 5. Priorização e Labels

Converta severidade em prioridade:

- Exposição de dados ou falha de segurança: `priority:critical`.
- Perda de dados ou falha funcional: `priority:high`.
- Degradação de comportamento: `priority:medium`.
- Dívida técnica sem sintoma ativo: `priority:low`.

Aplique também os labels de estado existentes no repositório (`status:ready` ou `status:needs-refinement`). Não crie labels novas sem necessidade.

Antes de criar cada issue, verifique duplicatas:

```bash
gh issue list --state open --search "<palavras-chave do achado>"
```

## 6. Criação no GitHub

Preserve a formatação usando um arquivo temporário:

```bash
gh issue create \
  --title "[Bug] Título" \
  --body-file /caminho/para/issue.md
```

## 7. Limites

- Não altere código durante a auditoria; a correção pertence ao executor.
- Não expanda escopo para refatorações abrangentes; registre como issue de dívida separada.
- Ao final, reporte de forma objetiva: issues criadas com número, achados que não viraram issue e o motivo, e áreas que não foram auditadas.
