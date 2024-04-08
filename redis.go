package dockerh

import (
	"fmt"
	"time"
)

// CreateRedis - creates a new redis pod
func CreateRedis(podName string, podPort int, password string) (ip string, err error) {

	Remove(podName)

	execParams := ""
	if len(password) > 0 {
		execParams = fmt.Sprintf("--requirepass %s", password)
	}

	err = Run(podName, "redis:latest", "", fmt.Sprintf("-d -p %d:6379", podPort), execParams)
	if err != nil {
		return
	}

	ip, err = WaitUntilListeningAndGetPodIP(podName, "", "", podPort, time.Minute, 3*time.Second)
	return
}
