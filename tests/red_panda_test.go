package dockerh_test

import (
	"testing"

	"github.com/rnojiri/dockerh"

	"github.com/stretchr/testify/assert"
)

// TestRedPanda - tests the red panda pod
func TestRedPanda(t *testing.T) {

	pod := "test-redpanda-pod"

	dockerh.Remove(pod)

	ip, err := dockerh.CreateRedPanda(pod, 9092)
	if !assert.NoError(t, err, "error starting redpanda pod") {
		return
	}

	defer dockerh.Remove(pod)

	validateIP(t, ip)
}
