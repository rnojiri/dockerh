package dockerh

import (
	"fmt"
	"strings"
	"time"
)

//
// RedPanda docker
// author: rnojiri
//

// CreateRedPanda - starts the red panda pod using default configurations
func CreateRedPanda(podName string, redPandaPort int) (string, error) {

	return CreateCustomRedPanda(podName, "", "", redPandaPort, defaultNoConnTimeout, defaultAfterConnTimeout)
}

// CreateRedPandaInNetwork - starts the red panda pod using default configurations
func CreateRedPandaInNetwork(podName, network string, redPandaPort int) (string, error) {

	return CreateCustomRedPanda(podName, "", network, redPandaPort, defaultNoConnTimeout, defaultAfterConnTimeout)
}

// CreateCustomRedPanda - starts the red panda pod using custom configurations
func CreateCustomRedPanda(podName, networkInspectFormat, network string, redPandaPort int, noConnTimeout, afterConnTimeout time.Duration) (string, error) {

	hostName := strings.ReplaceAll(strings.ToLower(podName), "-", "_")

	dockerArgs := fmt.Sprintf(
		"-p 8081:8081 -p 8082:8082 -p 9092:9092 -p 9644:9644 -v '%s:/var/lib/redpanda/data'",
		hostName,
	)

	execArgs := fmt.Sprintf(
		"redpanda start --smp 1  --memory 1G  --reserve-memory 0M --overprovisioned --set redpanda.empty_seed_starts_cluster=fals--seeds '%s:33145' --check=false --pandaproxy-addr INSIDE://0.0.0.0:28082,OUTSIDE://0.0.0.0:8082 --advertise-pandaproxy-addr INSIDE://%s:28082,OUTSIDE://localhost:8082 --kafka-addr INSIDE://0.0.0.0:29092,OUTSIDE://0.0.0.0:9092 --advertise-kafka-addr INSIDE://%s:29092,OUTSIDE://localhost:9092 --rpc-addr 0.0.0.0:33145 --advertise-rpc-addr %s:33145",
		hostName, hostName, hostName, hostName,
	)

	err := Run(podName, "docker.redpanda.com/vectorized/redpanda:latest", network, dockerArgs, execArgs)
	if err != nil {
		return "", err
	}

	return WaitUntilListeningAndGetPodIP(podName, networkInspectFormat, network, redPandaPort, noConnTimeout, afterConnTimeout)
}
