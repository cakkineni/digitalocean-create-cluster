package main

import (
	"strconv"
	//"net/http/httputil"
	//"flag"
	"code.google.com/p/goauth2/oauth"
	"fmt"
	"github.com/digitalocean/godo"
	"os"
	"time"
)

var (
	location,
	keyName,
	cloudConfigCluster,
	cloudConfigAgent,
	apiToken string
	serverCount int
	doClient    *godo.Client
	sshKeyId    []interface{}
)

func main() {

	login()

	cloudConfigCluster = createCloudConfigCluster()

	privateKey, publicKey := createSshKey()
	cloudConfigAgent = createCloudConfigAgent(publicKey)

	//create coreos servers
	var coreOSClusterDroplet *godo.DropletRoot
	for i := 0; i < serverCount; i++ {
		coreOSClusterDroplet = createCoreOSServer(i + 1)
	}

	//create agent server
	var pmxAgentDroplet *godo.DropletRoot
	pmxAgentDroplet = createAgentServer()

	println("Waiting for server creation")
	for {
		coreOSClusterDroplet, _, _ = doClient.Droplets.Get(coreOSClusterDroplet.Droplet.ID)
		if coreOSClusterDroplet.Droplet.Status == "active" {
			break
		}
		time.Sleep(60 * time.Millisecond)
	}

	println("Waiting for agent creation")
	for {
		pmxAgentDroplet, _, _ = doClient.Droplets.Get(pmxAgentDroplet.Droplet.ID)
		if pmxAgentDroplet.Droplet.Status == "active" {
			break
		}
		time.Sleep(60 * time.Millisecond)
	}


	agentIp := pmxAgentDroplet.Droplet.Networks.V4[1].IPAddress
	fleetIP := coreOSClusterDroplet.Droplet.Networks.V4[0].IPAddress

	setEtcdKey("agent-pri-ssh-key", privateKey)
	setEtcdKey("agent-fleet-api", agentIp)
	setEtcdKey("agent-public-ip", fleetIP)

	logout()

	//fmt.Scanln()
}

func init() {
	serverCount, _ = strconv.Atoi(os.Getenv("NODE_COUNT"))
	apiToken = os.Getenv("API_TOKEN")
	location = os.Getenv("REGION")
	keyName = os.Getenv("KEY_NAME")

	apiToken = "a37a4ba5a6ab6a9140bc2d1950776e901db71139fa59797ddd4deba57f5feabf"
	location = "nyc3"
	keyName = "macbook air"
	serverCount = 1

	if apiToken == "" || serverCount == 0 || location == "" {
		panic("\n\nMissing Params...Check Docs...\n\n")
	}
}

func login() {

	println("\nLogging in....")
	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: apiToken},
	}
	doClient = godo.NewClient(t.Client())

	intIds := []int{getSshKeyId()}
	for _, v := range intIds {
		sshKeyId = append(sshKeyId, v)
	}
}

func getSshKeyId() int {
	keys, _, err := doClient.Keys.List(&godo.ListOptions{Page: 1, PerPage: 10})
	keyId := -1

	if err != nil {
		panic(err)
	}

	for _, key := range keys {
		if key.Name == keyName {
			keyId = key.ID
			break
		}
	}

	if keyId == -1 {
		panic(fmt.Sprintf("\n\nSSH Key Name not found. Please make sure it matches exactly to your setup (case sensitive)\n\n"))
	}

	return keyId
}

func logout() {
	println("\nLogging out...Not really... TBD....")
}

func createCoreOSServer(id int) *godo.DropletRoot {
	println("Create CoreOS Server")
	var createReq *godo.DropletCreateRequest
	createReq = &godo.DropletCreateRequest{
		Name:              "coreos-" + strconv.Itoa(id),
		Region:            location,
		Size:              "512mb",
		Image:             "coreos-stable",
		PrivateNetworking: true,
		UserData:          cloudConfigCluster,
		SSHKeys:           sshKeyId,
	}
	return createServer(createReq)
}

func createAgentServer() *godo.DropletRoot {
	println("Create CoreOS Agent Server")
	var createReq *godo.DropletCreateRequest
	createReq = &godo.DropletCreateRequest{
		Name:              "pmx-remote-agent",
		Region:            location,
		Size:              "512mb",
		Image:             "coreos-stable",
		PrivateNetworking: true,
		UserData:          cloudConfigAgent,
		SSHKeys:           sshKeyId,
	}
	return createServer(createReq)
}

func createServer(createRequest *godo.DropletCreateRequest) *godo.DropletRoot {
	var err error
	newDroplet, _, err := doClient.Droplets.Create(createRequest)

	if err != nil {
		panic(err)
	}
	return newDroplet
}

func deleteServer(id int) {
	doClient.Droplets.Delete(id)
}

func newClient(token string) *godo.Client {
	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: token},
	}
	return godo.NewClient(t.Client())
}
