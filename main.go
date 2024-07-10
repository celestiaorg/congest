package main

import (
	"congest/cmd/netgen"
	"fmt"
	"log"
	"path/filepath"
	"sync"
	"time"

	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var (
	Regions = map[string]int{
		"nyc1": 3, "nyc3": 3, "sfo2": 3, "sfo3": 3, "ams3": 6, "sgp1": 7, "lon1": 4, "fra1": 3, "tor1": 3,
		"blr1": 7, "syd1": 7,
	}

	ReducedRegions = map[string]int{
		"nyc1": 1, "nyc3": 2, "sfo2": 2, "sfo3": 1, "ams3": 2, "sgp1": 3, "lon1": 2, "fra1": 2, "tor1": 2,
		"blr1": 2, "syd1": 2,
	}

	MinimalRegions = map[string]int{
		"nyc3": 1, "sfo3": 1, "ams3": 1, "sgp1": 1, "lon1": 1, "tor1": 1,
		"blr1": 1, "syd1": 1,
	}

	TestRegions = map[string]int{
		"sfo3": 1, "sgp1": 1,
	}

	sshIDs = []string{"31257644", "38506026", "32322687", "31138666", "22444021"}

	ChainID = "congest"
)

func main() {
	payloadRoot := "./payload"

	// Call the function to generate the network with the provided arguments
	n, err := netgen.NewNetwork(ChainID)
	if err != nil {
		panic(err)
	}

	wg := &sync.WaitGroup{}

	pulumi.Run(func(ctx *pulumi.Context) error {
		cursor := 0

		for region, count := range Regions {
			for j := 0; j < count; j++ {
				vname := fmt.Sprintf("validator-%d", cursor)
				ctx.Log.Info(fmt.Sprintf("Creating droplet %s in region %s", vname, region), nil)
				droplet, err := digitalocean.NewDroplet(ctx, vname, &digitalocean.DropletArgs{
					Region:  pulumi.String(region),
					Size:    pulumi.String("s-8vcpu-16gb"),     // Replace with the desired droplet size slug
					Image:   pulumi.String("ubuntu-22-04-x64"), // Replace with the desired image slug
					Name:    pulumi.String(vname),
					SshKeys: pulumi.ToStringArray(sshIDs),
				})
				if err != nil {
					ctx.Export("vname", pulumi.String(fmt.Sprintf("Error creating droplet %s %s: %s", vname, region, err.Error())))
					continue
				} else {
					ctx.Export(vname, droplet.Ipv4Address)
				}

				wg.Add(1)
				droplet.Ipv4Address.ApplyT(func(ip string) string {
					defer wg.Done()
					err = n.AddValidator(vname, ip, payloadRoot, region)
					if err != nil {
						panic(err)
					}

					return ip
				})

				cursor++
			}
		}
		go func() {
			wg.Wait()
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
		}()
		return nil
	})
	time.Sleep(10 * time.Second)

}
