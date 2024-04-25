package dockerh

import (
	"fmt"
	"time"
)

// CreateMemcached - creates a new memcached pod
func CreateMemcached(podName string, podPort, memoryMegabytes int) (ip string, err error) {

	Remove(podName)

	execParams := ""
	if memoryMegabytes > 0 {
		execParams = fmt.Sprintf("memcached -m %d", memoryMegabytes)
	}

	err = Run(podName, "memcached:latest", "", fmt.Sprintf("-d -p %d:11211", podPort), execParams)
	if err != nil {
		return
	}

	ip, err = WaitUntilListeningAndGetPodIP(podName, "", "", podPort, time.Minute, 3*time.Second)
	return
}
