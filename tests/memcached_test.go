package dockerh_test

import (
	"testing"

	"github.com/rnojiri/dockerh"

	"github.com/stretchr/testify/assert"
)

// TestMemcached - tests the memcached pod
func TestMemcached(t *testing.T) {

	pod := "test-memcached-pod"

	dockerh.Remove(pod)

	r, err := dockerh.CreateMemcached(pod, 11211, 64)
	if !assert.NoError(t, err, "error starting memcached pod") {
		return
	}

	dockerh.Remove(pod)

	validateIP(t, r)

	dockerh.Remove(pod)

	r, err = dockerh.CreateMemcached(pod, 11211, 0)
	if !assert.NoError(t, err, "error starting memcached pod with default memory") {
		return
	}

	dockerh.Remove(pod)

	validateIP(t, r)
}
