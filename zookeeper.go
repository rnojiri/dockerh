package dockerh

import (
	"fmt"
	"time"
)

//
// Zookeeper docker
// author: rnojiri
//

// CreateZookeeper - creates a zookeeper pod
func CreateZookeeper(podName string, port int) (string, error) {

	return CreateCustomZookeeper(podName, "", "", port, defaultNoConnTimeout, defaultAfterConnTimeout)
}

// CreateZookeeperInNetwork - starts the kafka pod using default configurations
func CreateZookeeperInNetwork(podName, network string, port int) (string, error) {

	return CreateCustomZookeeper(podName, "", network, port, defaultNoConnTimeout, defaultAfterConnTimeout)
}

// CreateCustomZookeeper - creates a zookeeper pod
func CreateCustomZookeeper(podName, networkInspectFormat, network string, port int, noConnTimeout, afterConnTimeout time.Duration) (string, error) {

	args := fmt.Sprintf("-e ZOOKEEPER_CLIENT_PORT=%d -e ZOOKEEPER_TICK_TIME=2000", port)

	err := Run(podName, "confluentinc/cp-zookeeper:latest", network, args, "")
	if err != nil {
		return "", err
	}

	return WaitUntilListeningAndGetPodIP(podName, networkInspectFormat, network, port, noConnTimeout, afterConnTimeout)
}

// CreateZookeeperWurstmeister - creates a zookeeper pod
func CreateZookeeperWurstmeister(podName string, port int) (string, error) {

	return CreateCustomZookeeperWurstmeister(podName, "", "", port, defaultNoConnTimeout, defaultAfterConnTimeout)
}

// CreateZookeeperWurstmeisterInNetwork - starts the kafka pod using default configurations
func CreateZookeeperWurstmeisterInNetwork(podName, network string, port int) (string, error) {

	return CreateCustomZookeeperWurstmeister(podName, "", network, port, defaultNoConnTimeout, defaultAfterConnTimeout)
}

// CreateCustomZookeeperWurstmeister - creates a zookeeper pod
func CreateCustomZookeeperWurstmeister(podName, networkInspectFormat, network string, port int, noConnTimeout, afterConnTimeout time.Duration) (string, error) {

	args := fmt.Sprintf("-p %d:2181", port)

	err := Run(podName, "wurstmeister/zookeeper:latest", network, args, "")
	if err != nil {
		return "", err
	}

	return WaitUntilListeningAndGetPodIP(podName, networkInspectFormat, network, port, noConnTimeout, afterConnTimeout)
}
