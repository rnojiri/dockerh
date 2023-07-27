package dockerh_test

import (
	"testing"

	"github.com/rnojiri/dockerh"

	"github.com/stretchr/testify/assert"
)

// TestPostgres - tests the postgres pod
func TestPostgres(t *testing.T) {

	pod := "test-postgres-pod"

	dockerh.Remove(pod)

	ip, err := dockerh.CreatePostgres(pod, "admin", "password", "db", 5432)
	if !assert.NoError(t, err, "error starting postgres pod") {
		return
	}

	defer dockerh.Remove(pod)

	validateIP(t, ip)
}
