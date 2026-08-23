# Skill: Refinamento de Ideias e Criação de Issues Ricas (QA)

Esta skill orienta o refinamento de novas ideias técnicas, análise de qualidade (QA) do repositório e a geração automatizada de GitHub Issues ricas e autossuficientes.

---

## 1. Fluxo de Refinamento de Ideias Brutais

Quando o usuário trouxer uma nova ideia, funcionalidade ou conceito bruto no chat, você deve atuar de forma **investigativa e provocativa**:

1.  **Faça perguntas inteligentes e desafiadoras**:
    *   Quais são os principais casos de borda (edge cases) e potenciais falhas dessa feature?
    *   Como essa alteração afeta a arquitetura ou performance das partes existentes?
    *   Há caminhos alternativos mais simples e eficientes para o mesmo objetivo?
    *   Quais tecnologias, dependências ou bibliotecas adicionais seriam introduzidas?
2.  **Provoque reflexões de viabilidade**: Não aceite de primeira sem analisar os impactos. Desafie premissas fracas.
3.  **Chegue a um consenso**: Só gere as issues após o usuário alinhar as respostas às suas provocações.

---

## 2. Fluxo de QA & Busca de Bugs

Ao analisar a qualidade do código para encontrar problemas (bugs, brechas de segurança, código ruim ou gargalos):
1.  **Audite os arquivos**: Use ferramentas de leitura do repositório para analisar os arquivos-fonte.
2.  **Categorize os problemas**: Identifique falhas específicas.
3.  **Gere Issues Atômicas**: Crie uma issue separada para cada problema ou grupo de problemas relacionados, evitando misturar escopos.

---

## 3. Template Estrito de Issue Rica

Toda issue criada por esta skill deve seguir exatamente esta estrutura Markdown para garantir que ela seja totalmente autossuficiente para o agente de escrita de código:

```markdown
# [Feature/Bug] Título Claro e Conciso

## 1. Contexto & Problema
[Explicação detalhada e inteligível do cenário, por que isso é necessário e qual o problema de negócio ou técnico resolvido]

## 2. Proposta de Solução
[Abordagem técnica recomendada para resolver a questão, descrevendo como deve funcionar]

## 3. Onde está Localizado
[Mapeamento exato de quais pastas, arquivos, classes, métodos ou linhas de código serão modificados/criados]

## 4. Passo a Passo da Implementação
- [ ] Passo 1...
- [ ] Passo 2...

## 5. Instrução de Autonomia (Importante)
> [!NOTE]
> Caso você precise de mais contexto técnico ou informações durante o desenvolvimento autônomo, busque ativamente no repositório antes de fazer alterações ou faça perguntas ao desenvolvedor principal no chat.
```

---

## 4. Comandos de Geração

Para automatizar a criação da issue refinada no GitHub usando a `gh` CLI, utilize o seguinte formato de comando:

```bash
gh issue create \
  --title "[Feature/Bug] <Título>" \
  --body "<Conteúdo da issue estruturada no template acima>"
```
