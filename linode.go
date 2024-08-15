package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/pulumi/pulumi-linode/sdk/v3/go/linode"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var _ Provider = &Linode{}

type Linode struct {
	sshIDs        []string
	globalTimeout pulumi.CustomTimeouts
}

// NewLinodeProvider creates a new Linode provider with the given SSH key IDs
// and timeout. The timeout should take the form of a string that can be parsed
// by pulumi, so "30m" would be 30 minutes, "1h" would be 1 hour, etc.
func NewLinodeProvider(timeout string) (*Linode, error) {
	rawLinodeSshIDs := os.Getenv("LINODE_SSH_KEY_IDS")
	if rawLinodeSshIDs == "" {
		fmt.Println("no raw ssh key")
		return nil, fmt.Errorf("No SSH IDs provided, please provide a list of SSH IDs in the LINODE_SSH_IDS environment variable")
	}

	fmt.Println("linode ssh key", rawLinodeSshIDs)

	sshIDs := strings.Split(rawLinodeSshIDs, ",")
	if len(sshIDs) == 0 {
		fmt.Println("no ssh keys from parsing")
		return nil, fmt.Errorf("No SSH IDs provided, please provide a list of SSH IDs in the LINODE_SSH_IDS environment variable")
	}

	fmt.Println("linode ssh key", sshIDs)

	return &Linode{
		sshIDs: sshIDs,
		globalTimeout: pulumi.CustomTimeouts{
			Delete: timeout,
		},
	}, nil
}

// CreateValidatorInstance creates a new validator instance in the given region.
// It fulfills the Provider interface. The returned string is the IP address of
// the new instance.
func (l *Linode) CreateValidatorInstance(ctx *pulumi.Context, name, region string) (pulumi.StringOutput, error) {
	instance, err := linode.NewInstance(ctx, name, &linode.InstanceArgs{
		Region:         pulumi.String(region),
		Type:           pulumi.String("g6-standard-8"),
		Image:          pulumi.String("linode/ubuntu22.04"),
		Label:          pulumi.String(name),
		AuthorizedKeys: pulumi.ToStringArray(l.sshIDs),
		Tags:           pulumi.StringArray{pulumi.String("temp")},
		StackscriptId:  pulumi.IntPtr(1439088),
		StackscriptData: pulumi.Map{
			"hostname": pulumi.String(name),
		},
	}, pulumi.Timeouts(&l.globalTimeout))

	if err != nil {
		ctx.Export("name", pulumi.String(fmt.Sprintf("Error creating instance %s %s: %s", name, region, err.Error())))
	} else {
		ctx.Export(name, instance.IpAddress)
	}

	return instance.IpAddress, nil
}

var (
	LinodeRegions = map[string]int{
		"ap-west":      1,
		"ca-central":   1,
		"ap-southeast": 1,
		"us-ord":       1,
		"fr-par":       1,
		"us-sea":       1,
		"br-gru":       1,
		"nl-ams":       1,
		"se-sto":       1,
		"es-mad":       1,
		"in-maa":       1,
		"jp-osa":       1,
		"it-mil":       1,
		"us-mia":       1,
		"id-cgk":       1,
		"us-lax":       1,
		"us-central":   1,
		"us-west":      1,
		"us-southeast": 1,
		"us-east":      1,
		"eu-west":      1,
		"ap-south":     1,
		"eu-central":   1,
		"ap-northeast": 1,
	}

	LinodeFullRegions = map[string]int{
		"ap-west":      1,
		"ca-central":   1,
		"ap-southeast": 1,
		"us-ord":       1,
		"us-sea":       2,
		"br-gru":       3,
		"se-sto":       2,
		"es-mad":       1,
		"in-maa":       2,
		"jp-osa":       3,
		"it-mil":       1,
		"us-mia":       1,
		"id-cgk":       2,
		"us-lax":       1,
		"us-central":   1,
		"us-southeast": 1,
		"ap-south":     2,
		"eu-central":   1,
		"ap-northeast": 1,
	}

	LinodeTestRegions = map[string]int{
		"us-central": 1,
		"eu-west":    1,
	}
)
