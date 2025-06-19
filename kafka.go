package dockerh

import (
	"fmt"
	"strings"
	"time"
)

//
// Kafka docker
// author: rnojiri
//

// CreateKafka - starts the kafka pod using default configurations
func CreateKafka(podName string, kafkaPort int, zookeeperHost string, zookeeperPort int) (string, error) {

	return CreateCustomKafka(podName, "", "", kafkaPort, zookeeperHost, zookeeperPort, defaultNoConnTimeout, defaultAfterConnTimeout)
}

// CreateKafkaInNetwork - starts the kafka pod using default configurations
func CreateKafkaInNetwork(podName, network string, kafkaPort int, zookeeperHost string, zookeeperPort int) (string, error) {

	return CreateCustomKafka(podName, "", network, kafkaPort, zookeeperHost, zookeeperPort, defaultNoConnTimeout, defaultAfterConnTimeout)
}

// CreateCustomKafka - starts the kafka pod using custom configurations
func CreateCustomKafka(podName, networkInspectFormat, network string, kafkaPort int, zookeeperHost string, zookeeperPort int, noConnTimeout, afterConnTimeout time.Duration) (string, error) {

	extraArgs := fmt.Sprintf(
		"-p %d:9092 -e KAFKA_BROKER_ID=1 -e KAFKA_ZOOKEEPER_CONNECT=%s:%d -e KAFKA_ADVERTISED_LISTENERS='PLAINTEXT://localhost:2%d,PLAINTEXT_HOST://localhost:%d' -e KAFKA_LISTENER_SECURITY_PROTOCOL_MAP='PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT' -e KAFKA_INTER_BROKER_LISTENER_NAME=PLAINTEXT -e KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1",
		kafkaPort, zookeeperHost, zookeeperPort, kafkaPort, kafkaPort,
	)

	err := Run(podName, "confluentinc/cp-kafka:latest", network, extraArgs, "")
	if err != nil {
		return "", err
	}

	return WaitUntilListeningAndGetPodIP(podName, networkInspectFormat, network, kafkaPort, noConnTimeout, afterConnTimeout)
}

type Topic struct {
	Name       string
	Partitions int
	Replicas   int
}

func (t Topic) String() string {

	return fmt.Sprintf("%s:%d:%d", t.Name, t.Partitions, t.Replicas)
}

// CreateKafkaWurstmeister - starts the kafka pod using default configurations
func CreateKafkaWurstmeister(podName string, kafkaPort int, zookeeperHost string, zookeeperPort int, topics ...Topic) (string, error) {

	return CreateCustomKafkaWurstmeister(podName, "", "", kafkaPort, zookeeperHost, zookeeperPort, defaultNoConnTimeout, defaultAfterConnTimeout, topics...)
}

// CreateKafkaWurstmeisterInNetwork - starts the kafka pod using default configurations
func CreateKafkaWurstmeisterInNetwork(podName, network string, kafkaPort int, zookeeperHost string, zookeeperPort int, topics ...Topic) (string, error) {

	return CreateCustomKafkaWurstmeister(podName, "", network, kafkaPort, zookeeperHost, zookeeperPort, defaultNoConnTimeout, defaultAfterConnTimeout, topics...)
}

// CreateCustomKafkaWurstmeister - starts the kafka pod using custom configurations
func CreateCustomKafkaWurstmeister(podName, networkInspectFormat, network string, kafkaPort int, zookeeperHost string, zookeeperPort int, noConnTimeout, afterConnTimeout time.Duration, topics ...Topic) (string, error) {

	kafkaTopics := make([]string, len(topics))

	if len(topics) > 0 {

		for i, item := range topics {

			kafkaTopics[i] = item.String()
		}
	}

	kafkaTopicsStr := ""

	if len(kafkaTopics) > 0 {

		kafkaTopicsStr = strings.Join(kafkaTopics, ",")
	}

	extraArgs := fmt.Sprintf(
		"-p %d:9092 -e KAFKA_ADVERTISED_HOST_NAME=kafka -e KAFKA_ADVERTISED_PORT=%d -e KAFKA_ZOOKEEPER_CONNECT=%s:%d -e KAFKA_CREATE_TOPICS=%s",
		kafkaPort, kafkaPort, zookeeperHost, zookeeperPort, kafkaTopicsStr,
	)

	err := Run(podName, "wurstmeister/kafka:latest", network, extraArgs, "")
	if err != nil {
		return "", err
	}

	return WaitUntilListeningAndGetPodIP(podName, networkInspectFormat, network, kafkaPort, noConnTimeout, afterConnTimeout)
}
