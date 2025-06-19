package dockerh_test

import (
	"testing"

	"github.com/rnojiri/dockerh"

	"github.com/stretchr/testify/assert"
)

// TestKafka - tests the kafka pod
func TestKafka(t *testing.T) {

	network := "test-kafka-network"
	podK := "test-kafka-pod"
	podZK := "test-zookeeper-pod"

	dockerh.Remove(podK)
	dockerh.Remove(podZK)
	dockerh.RemoveBridgeNetwork(network)

	err := dockerh.CreateBridgeNetwork(network)
	if !assert.NoError(t, err, "error starting network") {
		return
	}

	defer dockerh.RemoveBridgeNetwork(network)

	zkIP, err := dockerh.CreateZookeeperInNetwork(podZK, network, 2181)
	if !assert.NoError(t, err, "error starting zookeeper pod") {
		return
	}

	defer dockerh.Remove(podZK)

	validateIP(t, zkIP)

	ip, err := dockerh.CreateKafkaInNetwork(podK, network, 9092, podZK, 2181)
	if !assert.NoError(t, err, "error starting kafka pod") {
		return
	}

	defer dockerh.Remove(podK)

	validateIP(t, ip)
}

// TestKafkaWurstmeister - tests the kafka Wurstmeister pod
func TestKafkaWurstmeister(t *testing.T) {

	network := "test-kafka-wurstmeister-network"
	podK := "test-kafka-wurstmeister-pod"
	podZK := "test-zookeeper-pod"

	dockerh.Remove(podK)
	dockerh.Remove(podZK)
	dockerh.RemoveBridgeNetwork(network)

	err := dockerh.CreateBridgeNetwork(network)
	if !assert.NoError(t, err, "error starting network") {
		return
	}

	defer dockerh.RemoveBridgeNetwork(network)

	zkIP, err := dockerh.CreateZookeeperWurstmeisterInNetwork(podZK, network, 2181)
	if !assert.NoError(t, err, "error starting zookeeper pod") {
		return
	}

	defer dockerh.Remove(podZK)

	validateIP(t, zkIP)

	ip, err := dockerh.CreateKafkaWurstmeisterInNetwork(podK, network, 9092, podZK, 2181)
	if !assert.NoError(t, err, "error starting kafka wurstmeister pod") {
		return
	}

	defer dockerh.Remove(podK)

	validateIP(t, ip)
}
