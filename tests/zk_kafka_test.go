package dockerh_test

import (
	"context"
	"strings"
	"testing"
	"time"

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

// TestKafkaWurstmeisterReadyOnReturn - the broker and the topics must be ready as soon as the function returns (no sleeps)
func TestKafkaWurstmeisterReadyOnReturn(t *testing.T) {

	network := "test-kafka-ready-network"
	podK := "test-kafka-ready-pod"
	podZK := "test-zookeeper-ready-pod"
	topic := dockerh.Topic{Name: "ready-topic", Partitions: 3, Replicas: 1}

	dockerh.Remove(podK)
	dockerh.Remove(podZK)
	dockerh.RemoveBridgeNetwork(network)

	err := dockerh.CreateBridgeNetwork(network)
	if !assert.NoError(t, err, "error starting network") {
		return
	}

	defer dockerh.RemoveBridgeNetwork(network)

	_, err = dockerh.CreateZookeeperWurstmeisterInNetwork(podZK, network, 2181)
	if !assert.NoError(t, err, "error starting zookeeper pod") {
		return
	}

	defer dockerh.Remove(podZK)

	// zookeeper must answer right after returning
	out, err := dockerh.Exec(podZK, "echo ruok | nc localhost 2181")
	assert.NoError(t, err, "zookeeper command failed")
	assert.Contains(t, out, "imok", "expected zookeeper to be ok")

	_, err = dockerh.CreateKafkaWurstmeisterInNetwork(podK, network, 9092, podZK, 2181, topic)
	if !assert.NoError(t, err, "error starting kafka wurstmeister pod") {
		return
	}

	defer dockerh.Remove(podK)

	// the broker answers and the topic exists right after returning
	out, err = dockerh.Exec(podK, "kafka-topics.sh --bootstrap-server localhost:9092 --describe --topic "+topic.Name)
	if assert.NoError(t, err, "describe failed: %s", out) {
		assert.Contains(t, out, "PartitionCount: 3", "expected the topic partitions")
	}
}

// TestWaitUntilCommandSucceeds - tests the readiness command polling
func TestWaitUntilCommandSucceeds(t *testing.T) {

	pod := "test-ready-command-pod"

	rmPod(pod)

	err := dockerh.Run(pod, "alpine:latest", "", "", "sleep 120")
	if !assert.NoError(t, err, "error creating pod") {
		return
	}

	defer rmPod(pod)

	assert.NoError(t, dockerh.WaitUntilCommandSucceeds(pod, "echo imok", "imok", 5*time.Second))

	err = dockerh.WaitUntilCommandSucceeds(pod, "echo nope", "imok", 2*time.Second)
	assert.ErrorIs(t, err, dockerh.ErrPodNotReady)
	assert.True(t, strings.Contains(err.Error(), pod), "expected the pod in the error")
}

// TestWaitUntilCommandSucceedsHungCommand - a command that never returns must not outlive the total timeout
func TestWaitUntilCommandSucceedsHungCommand(t *testing.T) {

	pod := "test-ready-hung-command-pod"

	rmPod(pod)

	err := dockerh.Run(pod, "alpine:latest", "", "", "sleep 120")
	if !assert.NoError(t, err, "error creating pod") {
		return
	}

	defer rmPod(pod)

	timeout := 4 * time.Second
	start := time.Now()

	err = dockerh.WaitUntilCommandSucceeds(pod, "sleep 600", "", timeout)

	elapsed := time.Since(start)

	assert.ErrorIs(t, err, dockerh.ErrPodNotReady)
	assert.Less(t, elapsed, timeout+5*time.Second, "the hung command was not bounded by the timeout")
}

// TestExecWithTimeout - the docker exec client is killed when the timeout elapses
func TestExecWithTimeout(t *testing.T) {

	pod := "test-exec-timeout-pod"

	rmPod(pod)

	err := dockerh.Run(pod, "alpine:latest", "", "", "sleep 120")
	if !assert.NoError(t, err, "error creating pod") {
		return
	}

	defer rmPod(pod)

	out, err := dockerh.ExecWithTimeout(pod, "echo hi", 10*time.Second)
	assert.NoError(t, err)
	assert.Contains(t, out, "hi")

	start := time.Now()
	_, err = dockerh.ExecWithTimeout(pod, "sleep 600", 2*time.Second)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
	assert.Less(t, time.Since(start), 6*time.Second)
}
