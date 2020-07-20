package docker_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/uol/dockerh"

	"github.com/stretchr/testify/assert"
)

// TestScylla - tests the scylla pod
func TestScylla(t *testing.T) {

	pod := "test-scylla-pod"

	dockerh.Remove(pod)

	ip, err := dockerh.StartScylla(pod, "", "", 30*time.Second)
	if !assert.NoError(t, err, "error starting scylla pod") {
		return
	}

	defer dockerh.Remove(pod)

	assert.Regexp(t, regexp.MustCompile(`[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+`), ip, "expected some valid ip")
}
