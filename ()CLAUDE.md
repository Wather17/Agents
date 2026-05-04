# Work Workflow

## Session Start

List all files with `TASK` in the name by running `ls TASK-*.md` in alphabetical order — that is your execution queue.

Before anything else, read each `TASK-*.md` file — they define the full scope of what you will implement. If no task files exist or they are vague, log in `BLOCKERS.md` and wait for human input.

With the scope clear:

1. Run `git status` — analyze the repository state. If there are pending files, review the changes and commit before starting
2. Run the full test suite — only proceed if the repository is stable. If there are pre-existing failures, log them in `BLOCKERS.md` and move to the next task
3. Create a specific branch for the task — one task, one branch

## Execution

Break the task into smaller parts before writing any code. Think systematically: what needs to exist before what, what are the dependencies between parts.

Execute one part at a time, committing atomically along the way — do not save everything for the end.

### Tests

Every new implementation that can be tested must be tested. When writing tests consider:

- **Happy path** — the expected behavior in the normal flow
- **Edge cases** — empty, null, boundary, and unexpected inputs
- **Regression** — ensure what worked before still works
- **Integration** — if your implementation interacts with other modules, test the boundary

### Regression Analysis

Before wrapping up, actively map what may have been affected by your change:

- Functions or modules that call the code you changed
- Types and interfaces that depend on the modified structures
- Implicit behaviors that the previous code guaranteed
- Run the full test suite and analyze any failure — do not ignore warnings

## Escalation Criteria

During execution you will encounter two types of problems:

**You resolve on your own:**
- Syntax errors, linting, types
- Tests breaking due to the current implementation
- Missing dependency
- Refactor needed to complete the task within the defined scope

**You log in `BLOCKERS.md` and move to the next task:**
- Architectural decision not covered in the task
- Conflict with existing code that changes the original scope
- Missing credential, environment variable, or access
- Requirement ambiguity where any path has significant trade-offs

Blocker format:

```markdown
## Blocker — [timestamp]
**What you were doing:** description of the step in execution
**Problem found:** objective description
**Decision needed:** what needs to be decided or provided
**Affected files:** list of relevant files
```

## Continuity

You never sit idle. If a task hits a human-scope blocker:

1. Commit everything you have done so far on the current branch
2. Log the blocker in `BLOCKERS.md` using the standard format
3. Move to the next task in the queue

**Tasks are independent from each other.** If one task depends on the result of another, they should be grouped as a single task or as subtasks within the same entry. Dependency between separate tasks is a human planning mistake, not a problem for you to solve.

You only stop completely when:
- All tasks have been completed
- All remaining tasks are blocked waiting for human input

## Task Completion

With the implementation complete and all tests passing:

- Make sure your commits are atomic — related changes grouped, independent context separated
- Commit messages describe the **why**, not just the **what**
- Run the full test suite one last time before merging

## Merge

Merge the task branch into `develop`. After the merge:

```bash
git push origin develop
```

Clean up the task branch if it is no longer needed.

## Communication
Always respond in Portuguese (Brazilian), regardless of the language of these instructions.