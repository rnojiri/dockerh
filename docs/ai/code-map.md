# Code map

Where each feature lives, for agents looking for "the code that does X". Paths are repo-relative.

## Service helpers (public API)

| Feature | Entry points | Logic | Tests |
| --- | --- | --- | --- |
| Kafka (confluent cp-kafka) | `CreateKafka`, `CreateKafkaInNetwork`, `CreateCustomKafka` | `kafka.go` | `tests/zk_kafka_test.go` `TestKafka` (fails on master) |
| Kafka (wurstmeister) with topics | `CreateKafkaWurstmeister`, `...InNetwork`, `CreateCustomKafkaWurstmeister` | `kafka.go` (TCP wait, then broker and topic waits) | `TestKafkaWurstmeister`, `TestKafkaWurstmeisterReadyOnReturn` |
| Zookeeper (cp and wurstmeister) | `CreateZookeeper*`, `CreateZookeeperWurstmeister*` | `zookeeper.go` (wurstmeister also waits for `ruok`) | `tests/zk_kafka_test.go` |
| RedPanda | `CreateRedPanda*` | `redpanda.go` | `tests/red_panda_test.go` |
| Redis, MySQL, Postgres, Scylla, Memcached | `CreateRedis`, `CreateMysql*`, `CreatePostgres*`, `CreateScylla*`, `CreateMemcached` | `redis.go`, `mysql.go`, `postgres.go`, `scylla.go`, `memcached.go` | `tests/<service>_test.go` |

## Shared code

| Area | Path | Notes |
| --- | --- | --- |
| Run / remove / exists / networks | `docker.go` | `Run`, `Remove`, `Exists`, `CreateBridgeNetwork`, `RemoveBridgeNetwork` |
| TCP wait | `docker.go` | `WaitUntilListening`, `WaitUntilListeningAndGetPodIP`, `GetIPs`, `Address` |
| Command wait | `docker.go` | `Exec`, `WaitUntilCommandSucceeds`, `ErrPodNotReady`, `DefaultReadyTimeout` |
| Test helpers | `tests/docker_test.go` | `runPod`, `rmPod`, `validateIP`; tests for `docker.go` |

## Hotspots

| File | Why it changes often |
| --- | --- |
| `docker.go` | plumbing and wait logic shared by every helper |
| `kafka.go`, `zookeeper.go` | Wurstmeister support, topics, readiness |
| `go.mod` | dependency refreshes |
