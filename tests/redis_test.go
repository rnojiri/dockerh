package dockerh_test

import (
	"testing"

	"github.com/rnojiri/dockerh"

	"github.com/stretchr/testify/assert"
)

// TestRedis - tests the redis pod
func TestRedis(t *testing.T) {

	pod := "test-redis-pod"

	dockerh.Remove(pod)

	r, err := dockerh.CreateRedis(pod, 6379, "password")
	if !assert.NoError(t, err, "error starting redis pod") {
		return
	}

	defer dockerh.Remove(pod)

	validateIP(t, r)
}
