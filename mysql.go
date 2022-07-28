package dockerh

import (
	"fmt"
	"time"
)

//
// Mysql docker
// author: rnojiri
//

// CreateMysql - starts the mysql pod using custom configurations
func CreateMysql(podName, password, database string, port int) (string, error) {

	return CreateCustomMysql(podName, "", "", password, database, port, defaultWaitingTimeout)
}

// CreateMysqlInNetwork - starts the mysql pod using custom configurations
func CreateMysqlInNetwork(podName, network, password, database string, port int) (string, error) {

	return CreateCustomMysql(podName, "", network, password, database, port, defaultWaitingTimeout)
}

// CreateCustomMysql - starts the mysql pod
func CreateCustomMysql(podName, networkInspectFormat, network, password, database string, port int, timeout time.Duration) (string, error) {

	extraArgs := fmt.Sprintf("-p %d:%d -e MYSQL_ROOT_PASSWORD=%s -e MYSQL_DATABASE=%s", port, port, password, database)

	err := Run(podName, "mysql", network, extraArgs)
	if err != nil {
		return "", nil
	}

	return WaitUntilListeningAndGetPodIP(podName, networkInspectFormat, network, port, timeout)
}
