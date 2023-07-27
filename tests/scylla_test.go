package dockerh_test

import (
	"testing"

	"github.com/rnojiri/dockerh"

	"github.com/stretchr/testify/assert"
)

// TestScylla - tests the scylla pod
func TestScylla(t *testing.T) {

	pod := "test-scylla-pod"

	dockerh.Remove(pod)

	ip, err := dockerh.CreateScylla(pod, "")
	if !assert.NoError(t, err, "error starting scylla pod") {
		return
	}

	defer dockerh.Remove(pod)

	validateIP(t, ip)
}
