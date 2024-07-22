package main

import cmtconfig "github.com/tendermint/tendermint/config"

type ConfigOption func(*cmtconfig.Config)

type Experiment struct {
	CfgOptions []ConfigOption
	Regions    map[string]int
}

var (
	Experiments = map[string]Experiment{
		"100Nodes8MB": {
			Regions: DOFullRegions,
		},
		"MinimalNodes8MB": {
			Regions: DOMinimalRegions,
		},
	}
)

func GetExperiment(test string) (Experiment, bool) {
	experiment, ok := Experiments[test]
	return experiment, ok
}
