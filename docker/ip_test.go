package docker

import (
	"testing"

	"github.com/stvp/assert"
)

func TestRemoveIPSuffix(t *testing.T) {
	assert.Equal(t, "192.168.1.1", removeIPSuffix("192.168.1.1/24"))
	assert.Equal(t, "192.168.1.1", removeIPSuffix("192.168.1.1"))
}

func TestFindServiceIP(t *testing.T) {
	serviceName := []string{"TestFindServiceIP"}
	networkName := []string{"TestFindServiceIP"}
	network, _ := CreateNetwork(networkName)
	defer DeleteNetwork(networkName)
	StartService(&ServiceOptions{
		Namespace:  serviceName,
		Image:      "nginx",
		NetworksID: []string{network.ID},
	})
	defer StopService(serviceName)

	IP, err := FindServiceIP(networkName, serviceName)
	assert.Nil(t, err)
	assert.NotEqual(t, "", IP)
}

func TestFindServiceIPMissingNetwork(t *testing.T) {
	serviceName := []string{"TestFindServiceIPMissingNetwork"}
	networkName := []string{"TestFindServiceIPMissingNetwork"}
	StartService(&ServiceOptions{
		Namespace: serviceName,
		Image:     "nginx",
	})
	defer StopService(serviceName)

	IP, err := FindServiceIP(networkName, serviceName)
	assert.Nil(t, err)
	assert.Equal(t, "", IP)
}

func TestFindServiceIPMissingService(t *testing.T) {
	serviceName := []string{"TestFindServiceIPMissingService"}
	networkName := []string{"TestFindServiceIPMissingService"}
	CreateNetwork(networkName)
	defer DeleteNetwork(serviceName)

	IP, err := FindServiceIP(networkName, serviceName)
	assert.Nil(t, err)
	assert.Equal(t, "", IP)
}

func TestFindServiceIPWrongNetwork(t *testing.T) {
	serviceName := []string{"TestFindServiceIPWrongNetwork"}
	networkName := []string{"TestFindServiceIPWrongNetwork"}
	wrongNetworkName := []string{"TestFindServiceIPWrongNetwork", "DO NOT EXIST"}
	network, _ := CreateNetwork(networkName)
	defer DeleteNetwork(networkName)
	CreateNetwork(wrongNetworkName)
	defer DeleteNetwork(wrongNetworkName)
	StartService(&ServiceOptions{
		Namespace:  serviceName,
		Image:      "nginx",
		NetworksID: []string{network.ID},
	})
	defer StopService(serviceName)

	IP, err := FindServiceIP(wrongNetworkName, serviceName)
	assert.Nil(t, err)
	assert.Equal(t, "", IP)
}

func TestFindContainerIP(t *testing.T) {
	serviceName := []string{"TestFindContainerIP"}
	networkName := []string{"TestFindContainerIP"}
	network, _ := CreateNetwork(networkName)
	defer DeleteNetwork(networkName)
	StartService(&ServiceOptions{
		Namespace:  serviceName,
		Image:      "nginx",
		NetworksID: []string{network.ID},
	})
	defer StopService(serviceName)

	err := WaitForContainer(serviceName)
	assert.Nil(t, err)

	IP, err := FindContainerIP(networkName, serviceName)
	assert.Nil(t, err)
	assert.NotEqual(t, "", IP)
}
