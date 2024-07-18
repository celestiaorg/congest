package main

import (
	"congest/network"
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Provider interface {
	// CreateValidatorInstance creates a new validator instance in the given
	// region. The returned string is the IP address of the new instance.
	CreateValidatorInstance(ctx *pulumi.Context, name, region string) (pulumi.StringOutput, error)
}

// DeployValidators creates and deploys the given number of validators in the
// given regions. The cursor is the number of validators that have already been
// created. The function returns a map of the validator names to their IP
// addresses and the new cursor.
func DeployValidators(ctx *pulumi.Context, p Provider, regions map[string]int, cursor int) ([]network.NodeInfo, int) {
	validators := make([]network.NodeInfo, 0)
	for region, count := range regions {
		for j := 0; j < count; j++ {
			vname := fmt.Sprintf("validator-%d", cursor)

			ctx.Log.Info(fmt.Sprintf("Creating droplet %s in region %s", vname, region), nil)
			ip, err := p.CreateValidatorInstance(ctx, vname, region)
			if err != nil {
				ctx.Log.Error(fmt.Sprintf("Error creating droplet %s in region %s: %s", vname, region, err.Error()), nil)
				continue
			}
			validators = append(validators, network.NodeInfo{
				Name:      vname,
				PendingIP: ip,
				Region:    region,
			})

			cursor++
		}
	}

	return validators, cursor
}
