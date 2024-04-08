package dockerh_test

import (
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/rnojiri/dockerh"

	"github.com/stretchr/testify/assert"
)

var ipRegexp *regexp.Regexp = regexp.MustCompile(`[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+`)

//
// Has tests for the docker utility functions.
// author: rnojiri
//

func rmPod(pod string) {

	exec.Command("/bin/sh", "-c", "docker rm -f "+pod).Run()
}

func stopPod(pod string) {

	exec.Command("/bin/sh", "-c", "docker stop "+pod).Run()
}

func startPod(pod string) {

	exec.Command("/bin/sh", "-c", "docker start "+pod).Run()
}

func runPod(t *testing.T, pod, image string) bool {

	err := exec.Command("/bin/sh", "-c", fmt.Sprintf("docker run -d --name %s %s", pod, image)).Run()

	return assert.NoError(t, err, "error creating pod")
}

func psPod(t *testing.T, pod string) (string, bool) {

	output, err := exec.Command("/bin/sh", "-c", "docker ps -a -q --filter \"name="+pod+"\"").Output()

	if !assert.NoError(t, err, "error checking pod") {
		return "", false
	}

	return string(output), true
}

func podIP(t *testing.T, pod string) (string, bool) {

	output, err := exec.Command("/bin/sh", "-c", "docker inspect --format='{{ .NetworkSettings.Networks.bridge.IPAddress }}' "+pod).Output()

	if !assert.NoError(t, err, "error inspecting pod") {
		return "", false
	}

	return string(output), true
}

func networkLS(t *testing.T, network string) string {

	output, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("docker network ls | grep %s", network)).Output()

	if !assert.NoError(t, err, "error listing networks") {
		return ""
	}

	return string(output)
}

func podExists(t *testing.T, pod string) bool {

	output, ok := psPod(t, pod)
	if !ok {
		return false
	}

	return assert.Regexp(t, regexp.MustCompile("[0-9a-f]{12}"), output, "expected the pod's hash")
}

// TestRun - tests run command
func TestRun(t *testing.T) {

	pod := "test-run-hello-world"

	rmPod(pod)

	err := dockerh.Run(pod, "hello-world", "", "", "")
	if assert.NoError(t, err, "error not expected") {
		return
	}

	defer rmPod(pod)

	podExists(t, pod)
}

// TestRemove - tests remove command
func TestRemove(t *testing.T) {

	pod := "test-rm-hello-world"

	rmPod(pod)

	if !runPod(t, pod, "hello-world") {
		return
	}

	if !podExists(t, pod) {
		return
	}

	err := dockerh.Remove(pod)
	if !assert.NoError(t, err, "error not expected") {
		return
	}

	output, ok := psPod(t, pod)
	if !assert.True(t, ok, "error on ps") {
		return
	}

	assert.Equal(t, "", output, "expected no output")
}

// TestGetIPs - tests get ips command
func TestGetIPs(t *testing.T) {

	pod := "test-grafana"

	rmPod(pod)

	if !runPod(t, pod, "grafana/grafana") {
		return
	}

	defer rmPod(pod)

	if !podExists(t, pod) {
		return
	}

	format := ".NetworkSettings.Networks.bridge.IPAddress"

	libIps, err := dockerh.GetIPs(format, "", pod)
	if !assert.NoError(t, err, "error not expected", err) {
		return
	}

	if !assert.Len(t, libIps, 3, "expected 3 ips") {
		return
	}

	inspect, ok := podIP(t, pod)
	if !assert.True(t, ok, "error not expected") {
		return
	}

	assert.Contains(t, libIps, dockerh.NewAddress(strings.ReplaceAll(inspect, "\n", ""), 0, false), "expected same IP")
}

// TestExists - tests exists command
func TestExists(t *testing.T) {

	pod := "test-exists-http-https-echo"

	rmPod(pod)

	if !runPod(t, pod, "mendhak/http-https-echo") {
		return
	}

	if !podExists(t, pod) {
		return
	}

	running, err := dockerh.Exists(pod, dockerh.Running)
	if !assert.NoError(t, err, "error not expected") {
		return
	}

	if !assert.True(t, running, "expected pod running") {
		return
	}

	stopPod(pod)

	exited, err := dockerh.Exists(pod, dockerh.Exited)
	if !assert.NoError(t, err, "error not expected") {
		return
	}

	if !assert.True(t, exited, "expected pod exited") {
		return
	}

	startPod(pod)

	running, err = dockerh.Exists(pod, dockerh.Running)
	if !assert.NoError(t, err, "error not expected") {
		return
	}

	if !assert.True(t, running, "expected pod running") {
		return
	}

	rmPod(pod)

	notFound, err := dockerh.Exists(pod, dockerh.NotFound)
	if !assert.NoError(t, err, "error not expected") {
		return
	}

	if !assert.True(t, notFound, "expected pod not exists") {
		return
	}
}

// TestWaitUntilListening - wait until the host and port is listening
func TestWaitUntilListening(t *testing.T) {

	address := dockerh.NewAddress("localhost", 18123, false)

	go func() {
		<-time.After(1 * time.Second)

		listener, err := net.Listen("tcp", address.GetAddress())
		if assert.NoError(t, err, "expected no error when listening") {
			return
		}

		defer listener.Close()

		c, err := listener.Accept()
		if err != nil {
			if assert.NoError(t, err, "expected no error when accepting a connection") {
				return
			}
		}

		defer c.Close()

		<-time.After(1 * time.Second)
	}()

	connected := dockerh.WaitUntilListening(3*time.Second, 1*time.Second, address)

	assert.Len(t, connected, 1, "expected the address connection")
	assert.True(t, connected[0].Connected, "expected connected true")
	assert.False(t, connected[0].Fallback, "expected no fallback")
	assert.Equal(t, "localhost", connected[0].Host, "expected same host")
	assert.Equal(t, 18123, connected[0].Port, "expected same port")
}

func TestBridgeNetwork(t *testing.T) {

	expected := "dockerh-network-test"

	dockerh.RemoveBridgeNetwork(expected)

	err := dockerh.CreateBridgeNetwork(expected)
	if assert.NoError(t, err, "expected no error when creating a network") {
		return
	}

	res := networkLS(t, expected)
	if !assert.True(t, len(res) > 0, "expected some result") {
		return
	}

	assert.Equal(t, expected, res, "expected to find the network")

	err = dockerh.RemoveBridgeNetwork(expected)
	if assert.NoError(t, err, "expected no error when removing a network") {
		return
	}

	res = networkLS(t, expected)
	if !assert.True(t, len(res) == 0, "expected no results") {
		return
	}
}

func validateIP(t *testing.T, ip string) {

	assert.Regexp(t, ipRegexp, ip, "expected some valid ip")
}
