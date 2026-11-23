# Code review instructions

Per-layer conventions for this codebase, written for AI agents and humans. Agents read them because `AGENTS.md`
points here; the `applyTo` globs are matched by the agent. When several files apply, all are in force.

## Index

| File | Applies to |
| --- | --- |
| [`gympass-libraries.md`](./gympass-libraries.md) | `**` - no Gympass libraries are used here |
| [`helpers.md`](./helpers.md) | `*.go`, `tests/**` - service helpers and readiness |
