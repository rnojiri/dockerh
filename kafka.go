package dockerh

import (
	"fmt"
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
		"-p %d:%d -e KAFKA_BROKER_ID=1 -e KAFKA_ZOOKEEPER_CONNECT=%s:%d -e KAFKA_ADVERTISED_LISTENERS='PLAINTEXT://localhost:2%d,PLAINTEXT_HOST://localhost:%d' -e KAFKA_LISTENER_SECURITY_PROTOCOL_MAP='PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT' -e KAFKA_INTER_BROKER_LISTENER_NAME=PLAINTEXT -e KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1",
		kafkaPort, kafkaPort, zookeeperHost, zookeeperPort, kafkaPort, kafkaPort,
	)

	err := Run(podName, "confluentinc/cp-kafka:latest", network, extraArgs)
	if err != nil {
		return "", err
	}

	return WaitUntilListeningAndGetPodIP(podName, networkInspectFormat, network, kafkaPort, noConnTimeout, afterConnTimeout)
}
