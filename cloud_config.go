package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func createCloudConfigCluster() string {
	println("Create Cloud Config Cluster")
	response, _ := http.Get("https://discovery.etcd.io/new")
	defer response.Body.Close()
	contents, _ := ioutil.ReadAll(response.Body)
	cloudConfig, _ := ioutil.ReadFile("cloud-config-init.yaml")
	discoveryUrl := fmt.Sprintf("discovery: %s", string(contents))
	cloudConfigNew := strings.Replace(string(cloudConfig), "discovery_url", discoveryUrl, -1)
	return string(cloudConfigNew)
}

func createCloudConfigAgent() string {
	println("Create Cloud Config Agent")
	pubKey, _ := ioutil.ReadFile("~/.ssh/id_rsa.pub")
	cloudConfig, _ := ioutil.ReadFile("cloud-config-agent.yaml")
	cloudConfigNew := strings.Replace(string(cloudConfig), "ssh-rsa", fmt.Sprintf("ssh-rsa: %s", pubKey), -1)
	return string(cloudConfigNew)
}
