package main

import (
	"fmt"

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
func NewDigitalOcean(sshIDs []string, timeout string) *DigitalOcean {
	return &DigitalOcean{
		sshIDs: sshIDs,
		globalTimeout: pulumi.CustomTimeouts{
			Delete: timeout,
		},
	}
}

// CreateValidatorInstance creates a new validator instance in the given region.
// It fulfills the Provider interface. The returned string is the IP address of
// the new instance.
func (d *DigitalOcean) CreateValidatorInstance(ctx *pulumi.Context, name, region string) (pulumi.StringOutput, error) {
	droplet, err := digitalocean.NewDroplet(ctx, name, &digitalocean.DropletArgs{
		Region:  pulumi.String(region),
		Size:    pulumi.String("s-8vcpu-16gb"),     // Replace with the desired droplet size slug
		Image:   pulumi.String("ubuntu-22-04-x64"), // Replace with the desired image slug
		Name:    pulumi.String(name),
		SshKeys: pulumi.ToStringArray(sshIDs),
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
		"nyc1": 4, "nyc3": 4, "tor1": 4, "sfo2": 6, "sfo3": 6, "ams3": 11, "sgp1": 13, "lon1": 7, "fra1": 5,
		"blr1": 13, "syd1": 13,
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
		"nyc3": 1, "sfo3": 1, "ams3": 1, "sgp1": 1, "lon1": 1, "tor1": 1,
		"blr1": 1, "syd1": 1,
	}

	DOTestRegions = map[string]int{
		"sfo3": 1, "sgp1": 1,
	}
)

var (
	FullRegions = Regions{
		DO: DOFullRegions,
	}

	HalfRegions = Regions{
		DO: DOHalfRegions,
	}

	ReducedRegions = Regions{
		DO: DOReducedRegions,
	}

	MinimalRegions = Regions{
		DO: DOMinimalRegions,
	}

	TestRegions = Regions{
		DO: DOTestRegions,
	}
)
