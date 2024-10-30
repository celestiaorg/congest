package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var _ Provider = &DigitalOcean{}

type DigitalOcean struct {
	sshIDs        []string
	globalTimeout pulumi.CustomTimeouts
}

// NewDigitalOcean creates a new DigitalOcean provider with the given SSH IDs
// and timeout. The timeout should take the form of a string that can be parsed
// by pulumi, so "30m" would be 30 minutes, "1h" would be 1 hour, etc.
func NewDigitalOcean(timeout string) (*DigitalOcean, error) {
	rawDOSshIDs := os.Getenv("DO_SSH_KEY_IDS")
	if rawDOSshIDs == "" {
		fmt.Println("no raw ssh key")
		return nil, fmt.Errorf("No SSH IDs provided, please provide a list of SSH IDs in the DO_SSH_IDS environment variable")
	}

	sshIDs := strings.Split(rawDOSshIDs, ",")
	if len(sshIDs) == 0 {
		fmt.Println("no shh keys from parsing")
		return nil, fmt.Errorf("No SSH IDs provided, please provide a list of SSH IDs in the DO_SSH_IDS environment variable")
	}

	return &DigitalOcean{
		sshIDs: sshIDs,
		globalTimeout: pulumi.CustomTimeouts{
			Delete: timeout,
		},
	}, nil
}

// CreateValidatorInstance creates a new validator instance in the given region.
// It fulfills the Provider interface. The returned string is the IP address of
// the new instance.
func (d *DigitalOcean) CreateValidatorInstance(ctx *pulumi.Context, name, region string) (pulumi.StringOutput, error) {
	droplet, err := digitalocean.NewDroplet(ctx, name, &digitalocean.DropletArgs{
		Region:  pulumi.String(region),
		Size:    pulumi.String("c2-16vcpu-32gb"),   // Replace with the desired droplet size slug
		Image:   pulumi.String("ubuntu-22-04-x64"), // Replace with the desired image slug
		Name:    pulumi.String(name),
		SshKeys: pulumi.ToStringArray(d.sshIDs),
		Tags:    pulumi.ToStringArray([]string{"temp"}), // add a tag to make it easy to delete in the case that pulumi fails to delete the instance
	}, pulumi.Timeouts(&d.globalTimeout))
	if err != nil {
		ctx.Export("name", pulumi.String(fmt.Sprintf("Error creating droplet %s %s: %s", name, region, err.Error())))
	} else {
		ctx.Export(name, droplet.Ipv4Address)
	}

	return droplet.Ipv4Address, nil
}

var (
	DOFullRegions = map[string]int{
		"nyc3": 6, "tor1": 6, "sfo2": 6, "sfo3": 6, "ams3": 8, "sgp1": 2, "lon1": 8, "fra1": 8,
		"blr1": 2, "syd1": 2,
	}

	DOHalfRegions = map[string]int{
		"nyc1": 3, "nyc3": 3, "sfo2": 3, "sfo3": 3, "ams3": 6, "sgp1": 7, "lon1": 4, "fra1": 3, "tor1": 3,
		"blr1": 7, "syd1": 7,
	}

	DOReducedRegions = map[string]int{
		"nyc1": 1, "nyc3": 2, "sfo2": 2, "sfo3": 1, "ams3": 2, "sgp1": 3, "lon1": 2, "fra1": 2, "tor1": 2,
		"blr1": 2, "syd1": 2,
	}

	DOMinimalRegions = map[string]int{
		"nyc3": 2, "tor1": 1, "sfo3": 2, "ams3": 2, "lon1": 2, "fra1": 1,
	}

	DOTestRegions = map[string]int{
		"ams3": 1, "tor1": 1,
	}
)
