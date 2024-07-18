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
