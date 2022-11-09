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

	return CreateCustomPostgres(podName, "", "", user, password, database, port, defaultNoConnTimeout, defaultAfterConnTimeout)
}

// CreatePostgresInNetwork - starts the postgres pod using custom configurations
func CreatePostgresInNetwork(podName, network, user, password, database string, port int) (string, error) {

	return CreateCustomPostgres(podName, "", network, user, password, database, port, defaultNoConnTimeout, defaultAfterConnTimeout)
}

// CreateCustomPostgres - starts the postgres pod
func CreateCustomPostgres(podName, networkInspectFormat, network, user, password, database string, port int, noConnTimeout, afterConnTimeout time.Duration) (string, error) {

	extraArgs := fmt.Sprintf("-p %d:5432 -e POSTGRES_USER=%s -e POSTGRES_PASSWORD=%s -e POSTGRES_DB=%s", port, user, password, database)

	err := Run(podName, "postgres", network, extraArgs, "")
	if err != nil {
		return "", err
	}

	return WaitUntilListeningAndGetPodIP(podName, networkInspectFormat, network, port, noConnTimeout, afterConnTimeout)
}
