# AI docs change log

What the AI docs (`AGENTS.md`, `docs/`) learned from each batch of code changes, newest first. Maintained by the
`wh-repo-docs-and-graph-for-ai` skill: each entry covers the commits up to the one it names, and its `wh-repo-docs-and-graph-for-ai` comment is how
the skill knows which commits are already analysed. Keep those comments; when two branches add entries at the same
time, keep both entries.

## Log

### 2026-11-23 — up to `dd6213153a24`

<!-- wh-repo-docs-and-graph-for-ai: commit=dd6213153a244369783ea64e8f0c1f8b66ef0cde date=2026-11-23T16:23:02-03:00 fp=abb4e7259b1bfcac2bec66a1cb15fe10a6eeea8cb5bc87fcde0715d62b8c3655 -->

- First AI docs for the repo.
- fix(kafka): Wurstmeister Kafka waits for the broker (kafka-broker-api-versions.sh) and each topic (kafka-topics.sh --describe); Wurstmeister Zookeeper waits for ruok; new Exec, WaitUntilCommandSucceeds, ErrPodNotReady, DefaultReadyTimeout (2m).
- fix: WaitUntilListening applies afterConnTimeout with a single address.
- Known: TestKafka (cp-kafka) fails on master.

Docs updated: `AGENTS.md`, `CLAUDE.md`, `docs/ARCHITECTURE.md`, `docs/ai/code-map.md`, `docs/code-review-instructions/README.md`, `docs/code-review-instructions/gympass-libraries.md`, `docs/code-review-instructions/helpers.md`, `.gitignore`
