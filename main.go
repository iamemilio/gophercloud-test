package main

import (
	"fmt"
	"log"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/keypairs"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/utils/openstack/clientconfig"
)

func main() {
	opts := clientconfig.ClientOpts{
		Cloud: "optimist",
	}

	authOpts, err := clientconfig.AuthOptions(&opts)
	if err != nil {
		log.Fatalf("Auth Error: %v", err)
	}

	client, err := openstack.AuthenticatedClient(*authOpts)
	if err != nil {
		log.Fatalf("Client Error: %v", err)
	}

	serverClient, err := openstack.NewComputeV2(client, gophercloud.EndpointOpts{
		Region: "fra",
	})
	if err != nil {
		log.Fatalf("Server Client Error: %v", err)
	}

	testTags := []string{
		"test123",
	}

	serverOpts := servers.CreateOpts{
		Name:             "test",
		ImageName:        "Ubuntu 18.04 Bionic Beaver - Latest",
		FlavorRef:        "b7c4fa0b-7960-4311-a86b-507dbf58e8ac",
		AvailabilityZone: "es1",
		Tags:             testTags,
		Networks: []servers.Network{
			servers.Network{
				UUID: "5ef15f32-41ca-4dde-b541-e538199c28ca",
			},
		},
		SecurityGroups: []string{
			"kube-group",
		},
	}

	serverClient.Microversion = "2.52"
	server, err := servers.Create(serverClient, keypairs.CreateOptsExt{
		CreateOptsBuilder: serverCreateOpts,
		KeyName:           emilio,
	}).Extract()
	if err != nil {
		log.Fatalf("Server Create Error: %v", err)
	}

	err = servers.WaitForStatus(serverClient, server.ID, "ACTIVE", 240)
	if err != nil {
		log.Fatalf("Server Status Error: %v", err)
	}

	results := servers.Get(serverClient, server.ID)
	tags, err := results.ExtractTags()
	if err != nil {
		log.Fatalf("Tag Extraction Error: %v", err)
	}

	fmt.Printf("Tags: %s\n", tags)

	err = servers.Delete(serverClient, server.ID).ExtractErr()
	if err != nil {
		log.Fatalf("Tag Extraction Error: %v", err)
	}
}
