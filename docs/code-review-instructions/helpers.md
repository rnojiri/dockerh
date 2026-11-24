---
applyTo: "*.go,tests/**"
description: Service helpers, readiness checks and tests
---

# Service helpers

- Exported helper signatures MUST NOT change (library pinned by other repos). Add new behaviour behind the same
  signature, or a new function.
- A helper MUST NOT return before the service answers. For slow services (Kafka, Zookeeper) use
  `WaitUntilCommandSucceeds` after the TCP wait; never rely on `afterConnTimeout` as a readiness guess.
- A change for one service MUST NOT alter another service's helper.
- Command waits MUST fail with `ErrPodNotReady` (wrapped with `%w`) and name the pod, timeout and last output.
- Tests start real containers: remove leftovers first, `defer dockerh.Remove(...)`, no `time.Sleep` to wait for
  readiness (the helper is what is under test). Run with `go test -p 1`.
- `TestKafka` (cp-kafka) fails on `master`; verify against `master` before blaming a change.
- Run `gofmt -l .` and `go vet ./...` before pushing.
