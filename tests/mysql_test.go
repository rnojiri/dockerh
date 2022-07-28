package dockerh_test

import (
	"testing"

	"github.com/uol/dockerh"

	"github.com/stretchr/testify/assert"
)

// TestMysql - tests the mysql pod
func TestMysql(t *testing.T) {

	pod := "test-mysql-pod"

	dockerh.Remove(pod)

	ip, err := dockerh.CreateMysql(pod, "password", "db", 3306)
	if !assert.NoError(t, err, "error starting mysql pod") {
		return
	}

	defer dockerh.Remove(pod)

	validateIP(t, ip)
}
