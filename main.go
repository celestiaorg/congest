package main

import (
	"congest/network"
	"fmt"
	"log"
	"path/filepath"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	pcfg "github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

const (
	// TestConfigKey is the key used to retrieve which test is being ran from
	// the pulumi config. This value can be set by running `pulumi config set test TestName`
	TestConfigKey = "test"

	// RegionsConfgKey is the key used to retrieve the regions that the network
	// should be deployed to from the pulumi config. This value can be set by
	// running `pulumi config set regions Full`.
	RegionsConfgKey = "regions"

	// ChainIDConfigKey is the key used to retrieve the chain ID that the network
	// should be deployed with from the pulumi config. This value can be set by
	// running `pulumi config set chainID ChainID`.
	ChainIDConfigKey = "chainID"

	// GlobalTimeout is passed to all pulumi resources to ensure that they do
	// not stay alive too long.
	GlobalTimeoutString = "30m"
)

var (
	sshIDs = []string{"31257644", "38506026", "32322687", "31138666", "22444021"}

	ChainID = "congest"
)

func main() {
	payloadRoot := "./payload"

	// Call the function to generate the network with the provided arguments
	n, err := network.NewNetwork(ChainID)
	if err != nil {
		panic(err)
	}

	pulumi.Run(func(ctx *pulumi.Context) error {
		cursor := 0

		test := pcfg.Get(ctx, TestConfigKey)
		if test == "" {
			return fmt.Errorf("No test configuration provided, please assign a test configuration to the 'test' key in the Pulumi config with ")
		}

		do := NewDigitalOcean(sshIDs, GlobalTimeoutString)
		var validators []network.NodeInfo
		DOVals, cursor := DeployValidators(ctx, do, TestRegions.DO, cursor)
		validators = append(validators, DOVals...)

		ips := make([]pulumi.StringOutput, 0, len(validators))
		for _, val := range DOVals {
			n.AsyncAddValidator(val.Name, val.Region, payloadRoot, val.PendingIP)
			ips = append(ips, val.PendingIP)
		}
		pulumi.All(ips).ApplyT(func(_ []interface{}) error {
			err = n.InitNodes(payloadRoot)
			if err != nil {
				log.Fatal(err)
			}
			err = n.SaveValidatorsToFile(filepath.Join(payloadRoot, "validators.json"))
			if err != nil {
				log.Fatal(err)
			}
			err = n.SaveAddressBook(payloadRoot, n.Peers())
			if err != nil {
				log.Fatal(err)
			}
			return nil
		})
		return nil
	})

}

// Improvemnents - we can make better use of the apply txs I think by first
// doing all of the things that we can do then and there, then we wait to do all
// the sync stuff (like genesis creatiqon) there.

// We likely need to add multiple clounds

// we likely need to add things so that users can run specific tests instead of
// only being able to call pulumi up. I'd prefer just using go instead of config
// files for this as the mental overhead feels like less.

// as a random note: we need to get a mechanism that adds the ssh key to each of
// the nodes. With multiple clouds, I could see this being a bit of an issue.
// Perhaps this is where we use ansible? we jus ou ssh kay to on each node, and
// then we
