package dockerh

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

//
// Commons command execution functions.
// author: rnojiri
//

var (
	// ErrPodHashNotFound - raised when the output is incompatible with the hash pattern
	ErrPodHashNotFound error = errors.New("pod hash pattern was not found")

	// ErrNetworkHashNotFound - raised when the output is incompatible with the hash pattern
	ErrNetworkHashNotFound error = errors.New("network hash pattern was not found")

	// ErrPodIPNotFound - raised when the inspect format output  didn't returned any ip
	ErrPodIPNotFound error = errors.New("pod ip not found")

	// ErrPodNotListening - raised when the pod is not listening
	ErrPodNotListening error = errors.New("pod is not listening")

	regexpPodHashPattern         *regexp.Regexp = regexp.MustCompile("[a-f0-9]{64}")
	regexpDirtChars              *regexp.Regexp = regexp.MustCompile(`["'\r]+`)
	regexpDirtCharsAndLineBreaks *regexp.Regexp = regexp.MustCompile(`["'\r\n]+`)

	fallbackHosts []string = []string{"127.0.0.1", "0.0.0.0"}
)

// PodStatus - the pod status to be filtered
type PodStatus string

const (
	// Restarting - pod status
	Restarting PodStatus = "restarting"

	// Running - pod status
	Running PodStatus = "running"

	// Removing - pod status
	Removing PodStatus = "removing"

	// Paused - pod status
	Paused PodStatus = "paused"

	// Exited - pod status
	Exited PodStatus = "exited"

	// Dead - pod status
	Dead PodStatus = "dead"

	// DefaultnetworkInspectFormat - the default inspect format
	DefaultnetworkInspectFormat string = `(index .NetworkSettings.Networks "%s").IPAddress`

	defaultNetwork        string        = "bridge"
	defaultWaitingTimeout time.Duration = 1 * time.Minute
)

type Address struct {
	host string
	port int
}

func (a Address) GetAddress() string {

	return fmt.Sprintf("%s:%d", a.host, a.port)
}

// NewAddress - creates a new address
func NewAddress(host string, port int) Address {

	return Address{
		host: host,
		port: port,
	}
}

// createDockerCommand - creates a docker command to run or output
func createDockerCommand(cmd string) *exec.Cmd {

	return exec.Command("/bin/sh", "-c", fmt.Sprintf("docker %s", cmd))
}

// Run - runs a pod
func Run(name, image, network, extra string) error {

	networkParam := ""
	if len(network) > 0 {
		networkParam = fmt.Sprintf("--network %s", network)
	}

	output, err := createDockerCommand(fmt.Sprintf("run --name %s %s %s -d %s", name, networkParam, extra, image)).Output()
	if err != nil {
		return err
	}

	podHash := strings.Split(string(output), "\n")[0]

	if !regexpPodHashPattern.MatchString(podHash) {
		return ErrPodHashNotFound
	}

	return nil
}

// Remove - removes a pod
func Remove(pod string) error {

	return createDockerCommand(fmt.Sprintf("rm -f %s", pod)).Run()
}

// WaitUntilListeningAndGetPodIP - wait for the pod to be listening in the specied port and returns the pod's ip
func WaitUntilListeningAndGetPodIP(podName, networkInspectFormat, network string, port int, waitTimeout time.Duration) (string, error) {

	ips, err := GetIPs(networkInspectFormat, network, podName)
	if err != nil {
		return "", err
	}

	if len(ips) == 0 {
		return "", fmt.Errorf("%w: %s", ErrPodIPNotFound, podName)
	}

	addresses := make([]Address, len(ips))

	for i, ip := range ips {
		addresses[i] = NewAddress(ip, port)
	}

	connected := WaitUntilListening(waitTimeout, addresses...)

	if len(connected) == 0 {
		return "", fmt.Errorf("%w: %s:%d", ErrPodNotListening, podName, port)
	}

	return connected[0].host, nil
}

// GetIPs - return the pod's ips
func GetIPs(networkInspectFormat, network string, pod ...string) ([]string, error) {

	if len(network) == 0 {
		network = defaultNetwork
	}

	if len(networkInspectFormat) == 0 {
		networkInspectFormat = fmt.Sprintf(DefaultnetworkInspectFormat, network)
	}

	output, err := createDockerCommand(fmt.Sprintf("inspect --format='{{ %s }}' %s", networkInspectFormat, strings.Join(pod, " "))).Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(regexpDirtChars.ReplaceAllString(string(output), ""), "\n")

	all := []string{}
	all = append(all, fallbackHosts...)

	ips := lines[0 : len(lines)-1]
	for i, ip := range ips {
		if strings.Contains(ip, "<no value>") {
			ips = append(ips[:i], ips[i+1:]...)
		}
	}

	if len(ips) > 0 {
		all = append(all, ips...)
	}

	return all, nil
}

// Exists - check if a pod exists
func Exists(pod string, status PodStatus) (bool, error) {

	output, err := createDockerCommand(fmt.Sprintf(`ps -a -q --filter "name=%s" --filter "status=%s" --format "{{.Names}}"`, pod, status)).Output()
	if err != nil {
		return false, err
	}

	return regexpDirtCharsAndLineBreaks.ReplaceAllString(string(output), "") == pod, nil
}

// WaitUntilListening - wait some pod(s) to be listening
func WaitUntilListening(timeout time.Duration, addresses ...Address) (connected []Address) {

	if timeout == 0 {
		timeout = defaultWaitingTimeout
	}

	connected = []Address{}

	testConn := func(ctx context.Context, doneFunc context.CancelFunc, i int) {

		for {

			address := addresses[i].GetAddress()

			select {
			case <-ctx.Done():
				return
			default:
			}

			<-time.After(1 * time.Second)

			fmt.Println("trying: ", address)

			conn, err := net.DialTimeout("tcp", address, 1*time.Second)
			if err != nil {
				continue
			}

			if conn != nil {
				defer conn.Close()
				connected = append(connected, addresses[i])
				doneFunc()
				break
			}
		}
	}

	ctx, doneFunc := context.WithTimeout(context.Background(), timeout)

	for i := 0; i < len(addresses); i++ {

		go testConn(ctx, doneFunc, i)
	}

	<-ctx.Done()
	doneFunc = nil

	return connected
}

// WaitUntilListeningInSamePort - wait some pod(s) to be listening in a same port
func WaitUntilListeningInSamePort(timeout time.Duration, port int, hosts ...string) (connected []Address) {

	addresses := make([]Address, len(hosts))

	for i, h := range hosts {
		addresses[i] = Address{
			host: h,
			port: port,
		}
	}

	return WaitUntilListening(timeout, addresses...)
}

// CreateBridgeNetwork - creates a bridge network
func CreateBridgeNetwork(name string) error {

	output, err := createDockerCommand(fmt.Sprintf("network create %s", name)).Output()
	if err != nil {
		return err
	}

	networkHash := strings.Split(string(output), "\n")[0]

	if !regexpPodHashPattern.MatchString(networkHash) {
		return ErrNetworkHashNotFound
	}

	return nil
}

// RemoveBridgeNetwork - removes a bridge network
func RemoveBridgeNetwork(name string) error {

	_, err := createDockerCommand(fmt.Sprintf("network rm %s", name)).Output()

	return err
}
