package main

import (
	"congest/network"
	"log"
	"path/filepath"
	"time"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
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

func main() {
	payloadRoot := "./payload"

	experiment, chainID, err := readEnv()
	if err != nil {
		log.Fatal(err)
	}

	// Call the function to generate the network with the provided arguments
	n, err := network.NewNetwork(chainID)
	if err != nil {
		log.Fatal(err)
	}

	pulumi.Run(func(ctx *pulumi.Context) error {
		cursor := 0

		do, err := NewDigitalOcean(GlobalTimeoutString)
		if err != nil {
			log.Fatal(err)
		}

		linode, err := NewLinodeProvider(GlobalTimeoutString)
		if err != nil {
			log.Fatal(err)
		}

		var validators []network.NodeInfo

		DOVals, cursor := DeployValidators(ctx, do, experiment.Regions.DigitalOcean, cursor)
		validators = append(validators, DOVals...)

		linodeVals, cursor := DeployValidators(ctx, linode, experiment.Regions.Linode, cursor)
		validators = append(validators, linodeVals...)

		ips := make([]pulumi.StringOutput, 0, len(validators))

		for _, val := range validators {
			n.AsyncAddValidator(val.Name, val.Region, payloadRoot, val.PendingIP)
			ips = append(ips, val.PendingIP)
		}

		pulumi.All(ips).ApplyT(func(_ []interface{}) error {
			err = n.InitNodes(payloadRoot)
			if err != nil {
				log.Fatal(err)
			}

			err = n.SaveAddressBook(payloadRoot, n.Peers())
			if err != nil {
				log.Fatal(err)
			}

			err = n.SaveValidatorsToFile(filepath.Join(payloadRoot, "validators.json"))
			if err != nil {
				log.Fatal(err)
			}

			return nil
		})

		return nil
	})

	// Sleep for 30 seconds to allow for the instance to be accessible
	time.Sleep(30 * time.Second)
}
