package dockerh

import (
	"fmt"
	"time"
)

//
// Postgres docker
// author: rnojiri
//

// CreatePostgres - starts the postgres pod using custom configurations
func CreatePostgres(podName, user, password, database string, port int) (string, error) {

	return CreateCustomPostgres(podName, "", "", user, password, database, port, defaultWaitingTimeout)
}

// CreatePostgresInNetwork - starts the postgres pod using custom configurations
func CreatePostgresInNetwork(podName, network, user, password, database string, port int) (string, error) {

	return CreateCustomPostgres(podName, "", network, user, password, database, port, defaultWaitingTimeout)
}

// CreateCustomPostgres - starts the postgres pod
func CreateCustomPostgres(podName, networkInspectFormat, network, user, password, database string, port int, timeout time.Duration) (string, error) {

	extraArgs := fmt.Sprintf("-p %d:%d -e POSTGRES_USER=%s -e POSTGRES_PASSWORD=%s -e POSTGRES_DB=%s", port, port, user, password, database)

	err := Run(podName, "postgres", network, extraArgs)
	if err != nil {
		return "", nil
	}

	return WaitUntilListeningAndGetPodIP(podName, networkInspectFormat, network, port, timeout)
}
