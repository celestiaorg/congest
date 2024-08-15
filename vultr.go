package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dirien/pulumi-vultr/sdk/v2/go/vultr"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Vultr struct {
	sshIDs        []string
	globalTimeout pulumi.CustomTimeouts
}

// NewVultr creates a new Vultr provider with the given SSH IDs
// and timeout. The timeout should take the form of a string that can be parsed
// by pulumi, so "30m" would be 30 minutes, "1h" would be 1 hour, etc.
func NewVultr(timeout string) (*Vultr, error) {
	rawDOSshIDs := os.Getenv("VULTR_SSH_KEY_IDS")

	sshIDs := strings.Split(rawDOSshIDs, ",")
	if len(sshIDs) == 0 {
		return nil, fmt.Errorf("No SSH IDs provided, please provide a list of SSH IDs in the VULTR_SSH_KEY_IDS environment variable")
	}
	return &Vultr{
		sshIDs: sshIDs,
		globalTimeout: pulumi.CustomTimeouts{
			Delete: timeout,
		},
	}, nil
}

// CreateValidatorInstance creates a new validator instance in the given region.
// It fulfills the Provider interface. The returned string is the IP address of
// the new instance.
func (v *Vultr) CreateValidatorInstance(ctx *pulumi.Context, name, region string) (pulumi.StringOutput, error) {
	instance, err := vultr.NewInstance(ctx, name, &vultr.InstanceArgs{
		Region:    pulumi.String(region),
		Plan:      pulumi.String("vc2-6c-16gb"), // Replace with the desired instance plan
		OsId:      pulumi.Int(387),              // Replace with the desired OS ID (e.g., 387 for Ubuntu 22.04)
		Label:     pulumi.String(name),
		SshKeyIds: pulumi.ToStringArray(v.sshIDs),
		Hostname:  pulumi.String(name),
	}, pulumi.Timeouts(&v.globalTimeout))
	if err != nil {
		ctx.Export("name", pulumi.String(fmt.Sprintf("Error creating instance %s %s: %s", name, region, err.Error())))
	} else {
		ctx.Export(name, instance.MainIp)
	}

	return instance.MainIp, nil
}

var (
	VultrAllRegions = []string{
		"ams",
		"atl",
		"blr",
		"bom",
		"cdg",
		"del",
		"dfw",
		"ewr",
		"fra",
		"hnl",
		"icn",
		"itm",
		"jnb",
		"lax",
		"lhr",
		"mad",
		"man",
		"mel",
		"mex",
		"mia",
		"nrt",
		"ord",
		"sao",
		"scl",
		"sea",
		"sgp",
		"sjc",
		"sto",
		"syd",
		"tlv",
		"waw",
		"yto",
	}

	// VultrSlugs = map[string]string{
	// 	"ams": "vc2-6c-16gb",
	// 	"atl": "vc2-6c-16gb",
	// 	"blr",
	// 	"bom",
	// 	"cdg",
	// 	"del",
	// 	"dfw",
	// 	"ewr",
	// 	"fra",
	// 	"hnl",
	// 	"icn",
	// 	"itm",
	// 	"jnb",
	// 	"lax",
	// 	"lhr",
	// 	"mad",
	// 	"man",
	// 	"mel",
	// 	"mex",
	// 	"mia",
	// 	"nrt",
	// 	"ord",
	// 	"sao",
	// 	"scl",
	// 	"sea",
	// 	"sgp",
	// 	"sjc",
	// 	"sto",
	// 	"syd",
	// 	"tlv",
	// 	"waw",
	// 	"yto",
	// }

	VultrFullRegions = map[string]int{
		"ams": 1,
		"atl": 1,
		"blr": 1,
		"bom": 1,
		"cdg": 1,
		"del": 1,
		"dfw": 1,
		"fra": 1,
		// "hnl": 1, // doesn't support normal slug
		"icn": 1,
		"itm": 1,
		"jnb": 1,
		"lax": 1,
		"lhr": 1,
		"mad": 1,
		"man": 1,
		"mel": 1,
		"mex": 1,
		"mia": 1,
		"nrt": 1,
		"ord": 1,
		// "sao": 1, // doesn't support normal slug
		"scl": 1,
		"sea": 1,
		"sgp": 1,
		"sjc": 1,
		"sto": 1,
		"syd": 1,
		"tlv": 1,
		"yto": 1,
	}

	VultrHalfRegions = map[string]int{
		// "ams": 1,
		// "atl": 1,
		// "blr": 1,
		// "bom": 1,
		// "cdg": 1,
		// "del": 1,
		// "dfw": 1,
		"ewr": 1,
		"icn": 1,
		"itm": 1,
		"jnb": 1,
		"lax": 1,
		// "lhr": 1,
		"mad": 1,
		"man": 1,
		"mel": 1,
		"mex": 1,
		// "mia": 1,
		// "nrt": 1,
		// "ord": 1,
		// "sao": 1,
		// "scl": 1,
		// "sea": 1,
		// "sgp": 1,
		// "sjc": 1,
		// "sto": 1,
		// "syd": 1,
		// "tlv": 1,
		// "yto": 1,
	}

	VultrReducedRegions = map[string]int{
		"ams": 1,
		"atl": 1,
		"blr": 1,
		"bom": 1,
		"dfw": 1,
		"ewr": 1,
		"hnl": 1,
		"itm": 1,
		"jnb": 1,
		"lax": 1,
		"mex": 1,
		"mia": 1,
		"nrt": 1,
		"sao": 1,
		"scl": 1,
		"sea": 1,
		"sjc": 1,
	}

	VultrMinimalRegions = map[string]int{
		"blr": 1,
		"bom": 1,
		"cdg": 1,
		"del": 1,
		"dfw": 1,
		"ewr": 1,
		"nrt": 1,
		"ord": 1,
		"sao": 1,
		"scl": 1,
	}

	VultrTestRegions = map[string]int{
		"dfw": 1,
		"ord": 1,
		"scl": 1,
	}
)
