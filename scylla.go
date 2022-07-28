package dockerh

import (
	"time"
)

//
// Scylla docker
// author: rnojiri
//

// CreateScylla - starts the scylla pod using custom parameters
func CreateScylla(podName, extraCommands string) (string, error) {

	return CreateCustomScylla(podName, "", "", extraCommands, defaultWaitingTimeout)
}

// CreateScyllaInNetwork - starts the scylla pod using custom parameters
func CreateScyllaInNetwork(podName, network, extraCommands string) (string, error) {

	return CreateCustomScylla(podName, "", network, extraCommands, defaultWaitingTimeout)
}

// CreateCustomScylla - starts the scylla pod
func CreateCustomScylla(podName, networkInspectFormat, network, extraCommands string, timeout time.Duration) (string, error) {

	err := Run(podName, "scylladb/scylla", network, extraCommands)
	if err != nil {
		return "", nil
	}

	return WaitUntilListeningAndGetPodIP(podName, networkInspectFormat, network, 9042, timeout)
}
