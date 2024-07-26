package main

import (
	"fmt"
	"os"

	cmtconfig "github.com/tendermint/tendermint/config"
)

type ConfigOption func(*cmtconfig.Config)

type Experiment struct {
	CfgOptions []ConfigOption
	Regions    Regions
}

var (
	Experiments = map[string]Experiment{
		"100Nodes8MB": {
			Regions: FullRegions,
		},
		"MinimalNodes8MB": {
			Regions: MinimalRegions,
		},
		"Test8MB": {
			Regions: TestRegions,
		},
	}
)

func getExperiment(test string) (Experiment, bool) {
	experiment, ok := Experiments[test]
	return experiment, ok
}
func readEnv() (experiment Experiment, chainID string, err error) {
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
