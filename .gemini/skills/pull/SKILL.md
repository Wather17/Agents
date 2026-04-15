---
name: pull
description: Baixa as issues do GitHub via CLI e gera o artefato PLANO_EXECUCAO.md na raiz do projeto com checklists atômicos.
---

Precisamos resolver as issues pendentes do projeto. O terminal tem limite de histórico e pode esquecer instruções longas, então documentaremos tudo em arquivos físicos.

**Regras de Execução:**
1. Use o GitHub CLI para listar as issues abertas e exporte o resultado formatado para um arquivo usando:
`gh issue list --state open --json number,title,body --template '{{range .}}## #{{.number}} - {{.title}}{{"\n"}}{{.body}}{{"\n\n"}}{{end}}' > PENDENCIAS.md`
2. Leia o conteúdo de `PENDENCIAS.md` silenciosamente.
3. Com base na leitura, crie um novo arquivo chamado `PLANO_EXECUCAO.md` na raiz do projeto. 
4. Dentro do `PLANO_EXECUCAO.md`, monte um plano detalhado usando checklists de Markdown (`- [ ]`), priorizando os bugs críticos primeiro. Divida a resolução de cada issue em passos atômicos.
5. Quando terminar, me avise no chat e pergunte qual item do `PLANO_EXECUCAO.md` devemos atacar primeiro.
