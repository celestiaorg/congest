package main

import (
	"congest/cmd/netgen"
	"fmt"
	"path/filepath"

	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var (
	Regions = map[string]int{
		"nyc1": 2, "nyc3": 2, "sfo2": 2, "sfo3": 2, "ams3": 4, "sgp1": 5, "lon1": 4, "fra1": 4, "tor1": 3,
		"blr1": 6, "syd1": 6,
	}

	TestRegions = map[string]int{
		"sfo3": 1, "sgp1": 1,
	}

	sshIDs = []string{"31257644", "38506026", "32322687", "31138666", "22444021"}

	ChainID = "congest"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cursor := 0
		payloadRoot := "./payload"

		// Call the function to generate the network with the provided arguments
		n, err := netgen.NewNetwork(ChainID)
		if err != nil {
			return err
		}

		for region, count := range TestRegions {
			for j := 0; j < count; j++ {
				vname := fmt.Sprintf("validator-%d", cursor)
				droplet, err := digitalocean.NewDroplet(ctx, vname, &digitalocean.DropletArgs{
					Region:  pulumi.String(region),
					Size:    pulumi.String("s-8vcpu-16gb-intel"), // Replace with the desired droplet size slug
					Image:   pulumi.String("ubuntu-22-04-x64"),   // Replace with the desired image slug
					Name:    pulumi.String(vname),
					SshKeys: pulumi.ToStringArray(sshIDs),
				})
				if err != nil {
					return err
				}
				// Add outputs to get the droplet IP addresses
				ctx.Export(vname, droplet.Ipv4Address)

				var parsedIPV4 string
				droplet.Ipv4Address.ApplyT(func(ip string) string {
					parsedIPV4 = ip
					return ip
				})

				err = n.AddValidator(vname, parsedIPV4, payloadRoot)
				if err != nil {
					return err
				}

				cursor++
			}
		}

		n.InitNodes(payloadRoot)
		n.SaveValidatorsToFile(filepath.Join(payloadRoot, "validators.json"))
		err = n.SaveAddressBook(payloadRoot, n.Peers())
		if err != nil {
			return err
		}

		return nil
	})
}
