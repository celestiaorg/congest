package main

import (
	"congest/network"
	"fmt"
	"os"
	"time"

	cmtconfig "github.com/tendermint/tendermint/config"
)

var (
	Experiments = map[string]network.Experiment{
		"100Nodes8MB": {
			Regions: FullRegions,
		},
		"2MB6s": {
			Regions: FullRegions,
			CfgOptions: []network.ConfigOption{
				func(c *cmtconfig.Config) {
					// note!: these aren't actually used yet, but this is what they should look like imo
					c.Consensus.TimeoutCommit = time.Second * 4
					c.Consensus.TimeoutPropose = time.Second * 3
				},
			},
		},
		"HalfNodes8MB": {
			Regions: HalfRegions,
		},
		"MinimalNodes8MB": {
			Regions: MinimalRegions,
		},
		"Test8MB": {
			Regions: TestRegions,
		},
	}
)

func getExperiment(test string) (network.Experiment, bool) {
	experiment, ok := Experiments[test]
	return experiment, ok
}

func readEnv() (experiment network.Experiment, chainID string, err error) {
	chainID = os.Getenv("EXPERIMENT_CHAIN_ID")
	rawExperiment := os.Getenv("EXPERIMENT_NAME")

	if chainID == "" {
		return experiment, "", fmt.Errorf("No chain ID provided, please provide a chain ID in the EXPERIMENT_CHAIN_ID environment variable")
	}

	experiment, has := getExperiment(rawExperiment)
	if !has {
		return experiment, "", fmt.Errorf("No experiment found with the name %s", rawExperiment)
	}

	return experiment, chainID, nil
}