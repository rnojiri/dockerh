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

	return CreateCustomZookeeper(podName, "", "", port, defaultWaitingTimeout)
}

// CreateZookeeperInNetwork - starts the kafka pod using default configurations
func CreateZookeeperInNetwork(podName, network string, port int) (string, error) {

	return CreateCustomZookeeper(podName, "", network, port, defaultWaitingTimeout)
}

// CreateCustomZookeeper - creates a zookeeper pod
func CreateCustomZookeeper(podName, networkInspectFormat, network string, port int, timeout time.Duration) (string, error) {

	args := fmt.Sprintf("-e ZOOKEEPER_CLIENT_PORT=%d -e ZOOKEEPER_TICK_TIME=2000", port)

	err := Run(podName, "confluentinc/cp-zookeeper", network, args)
	if err != nil {
		return "", nil
	}

	return WaitUntilListeningAndGetPodIP(podName, networkInspectFormat, network, port, defaultWaitingTimeout)
}
