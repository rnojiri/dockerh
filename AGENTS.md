# dockerh

This file orients AI agents (Claude Code, Copilot, Cursor, Codex, ...) working in this repository. `CLAUDE.md` only
imports it, so every tool reads the same rules. Keep it short: details live in `docs/`.

`github.com/rnojiri/dockerh` is a small, dependency-light Go library (personal project of rnojiri) that starts
throwaway Docker containers (Kafka, Zookeeper, Redis, MySQL, Postgres, Scylla, Memcached, RedPanda) for integration
tests of other projects. It shells out to the `docker` CLI and waits until the container is usable. It has no
services, no datastore and no Gympass dependencies; it is consumed by Gympass repos such as identity-commons.
Releases are plain git tags (`v1.3.1`, `v1.4.0`, ...), consumers bump `go.mod`.

## Where to find what

| What you need | Where to look |
| --- | --- |
| Feature -> code (helpers, readiness, tests) | [`docs/ai/code-map.md`](./docs/ai/code-map.md) |
| Architecture and the readiness flow | [`docs/ARCHITECTURE.md`](./docs/ARCHITECTURE.md) |
| Review rules | [`docs/code-review-instructions/`](./docs/code-review-instructions/) |
| What changed recently | [`docs/ai/changes.md`](./docs/ai/changes.md) |
| Docker plumbing (`Run`, `Exec`, `WaitUntilListening`, `WaitUntilCommandSucceeds`, networks) | `docker.go` |
| One file per service helper | `kafka.go`, `zookeeper.go`, `redis.go`, `mysql.go`, `postgres.go`, `scylla.go`, `memcached.go`, `redpanda.go` |
| Tests (need a local Docker daemon) | `tests/`, run by `run-tests.sh` |

## Golden rules

- **Helpers are backward compatible**: the exported `Create*`, `Create*InNetwork` and `CreateCustom*` signatures are
  public API of a library other repos pin. Add behaviour, never change a signature (git history: tags, v1.4.0).
- **Return only when the service answers.** A TCP connect to the published port is not readiness: Docker's port proxy
  accepts connections before the process in the container started. Wurstmeister Kafka and Zookeeper therefore also run a
  readiness command inside the container (see ARCHITECTURE). A new helper for a slow-starting service does the same.
- **Stay lightweight**: no cgo and no heavy client libraries (e.g. no confluent-kafka-go); talk to services through
  `docker exec` and the standard library.
- **Do not change another service's helper** as a side effect of a fix for one service.

## Gympass libraries first

This repo imports no Gympass library (`go.mod` only has `testify`). Do not add one; the library must stay usable
by any project.

## Code style

Match the formatting and idioms of the surrounding code.

- Formatter: `gofmt -l .` must print nothing; `go vet ./...` must pass. There is no golangci-lint config or CI.
- Style visible in the code: a blank line after each function signature, `// Name - description` doc comments, one
  file per service, `Create<Service>` / `Create<Service>InNetwork` / `CreateCustom<Service>` triplets that delegate
  to the custom one (example: `zookeeper.go`), errors returned as is or wrapped with `%w` around a package sentinel
  (`ErrPodNotListening`, `ErrPodNotReady` in `docker.go`).

## Commands

| Command | Effect |
| --- | --- |
| `./run-tests.sh` | `go test -v -p 1 -count 1 -timeout 2m ./tests` (sequential, needs Docker; 2m is short for Kafka tests, raise it) |
| `go test -p 1 -count 1 -timeout 8m ./tests -run 'Kafka\|WaitUntil'` | Kafka and wait helpers only |
| `go vet ./... && gofmt -l .` | static checks |

Known failure: `TestKafka` (confluent cp-kafka helpers) fails on `master` too, independent of the Wurstmeister
readiness work; do not treat it as a regression.

## How AI-assisted work runs here

- **Before implementing:** read the matching `docs/code-review-instructions/` file.
- **Before pushing:** the `wh-repo-docs-and-graph-for-ai` skill updates this file and `docs/` from the new commits and records them
  in `docs/ai/changes.md` (a separate `docs: update AI guidelines` commit on the same branch).
- **Code graph:** when CodeGraph is installed (`codegraph version`), each clone keeps its own graph in
  `.codegraph/` (gitignored, never pushed). Prefer `codegraph explore` / its MCP tools to find code.
- Repo-local rules here outrank generic harness or user-level rules.
- Default branch: `master`. Never commit to it directly; open a PR from a branch. No PR template or labels here.
