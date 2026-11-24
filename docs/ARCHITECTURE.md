# Architecture

`dockerh` is a single Go package (`package dockerh`, repo root) plus `tests/` (external test package
`dockerh_test`). There is no daemon, store or config: every function runs `docker ...` through `/bin/sh -c`
(`createDockerCommand` in `docker.go`) or `docker exec` (`Exec`).

## Components

| Component | File | Role |
| --- | --- | --- |
| Docker plumbing | `docker.go` | `Run`, `Remove`, `Exists`, networks, `GetIPs`, TCP wait, command-based wait |
| Service helpers | `kafka.go`, `zookeeper.go`, `redis.go`, `mysql.go`, `postgres.go`, `scylla.go`, `memcached.go`, `redpanda.go` | image + docker args per service, then wait |

## Startup flow of a helper

1. `Run` starts the container detached (`docker run -d`), optionally in a network, and checks the container hash.
2. `WaitUntilListeningAndGetPodIP` collects candidate addresses (`GetIPs`: container IP, then fallbacks `127.0.0.1`
   and `0.0.0.0`) and calls `WaitUntilListening`, which dials each address once a second until `noConnTimeout`.
   After one address connects it pauses `afterConnTimeout` (default 3s) so the others can connect; the last address to
   connect does not pause, except when there is a single address (fixed: it used to skip the pause).
3. A TCP connect only proves Docker's port proxy accepts connections. Services that start slowly are then checked
   from inside the container.

## Readiness checks (Wurstmeister Kafka and Zookeeper)

`WaitUntilCommandSucceeds(pod, command, expectedOutput, timeout)` polls `Exec(pod, command)` (`docker exec pod sh -c`)
every second until it exits 0 and the output contains `expectedOutput`; on timeout it returns `ErrPodNotReady` with
the last error and output. `timeout <= 0` means `DefaultReadyTimeout` (exported var, 2 minutes, change it to tune).

| Helper | Check | Why |
| --- | --- | --- |
| `CreateCustomKafkaWurstmeister` (and the non-custom/`InNetwork` wrappers) | `kafka-broker-api-versions.sh --bootstrap-server localhost:9092` | the broker answers (CI evidence: kafkap tests failed with `1/1 brokers are down` after 3s because Kafka needed ~12s) |
| same, per `Topic` passed | `kafka-topics.sh --bootstrap-server localhost:9092 --describe --topic <name>` | the image creates `KAFKA_CREATE_TOPICS` topics in background after the broker starts |
| `CreateCustomZookeeperWurstmeister` | `echo ruok \| nc localhost 2181` expects `imok` | zookeeper is up (the `nc` binary exists in the image) |

Not covered: the confluent `CreateCustomKafka` / `CreateCustomZookeeper` helpers (cp images, different ports and
tools) and the other services (they start fast enough with the TCP wait). Their `TestKafka` test currently fails on
`master` as well.

## Tests

`tests/*_test.go` start real containers (Docker daemon required, `go test -p 1` so ports do not clash). Each test
removes leftovers first and defers `dockerh.Remove`.
