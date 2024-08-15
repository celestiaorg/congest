package main

import "congest/network"

var (
	FullRegions = network.Regions{
		DigitalOcean: DOFullRegions,
		Linode:       LinodeFullRegions,
	}

	HalfRegions = network.Regions{
		DigitalOcean: DOHalfRegions,
		Vultr:        VultrHalfRegions,
	}

	ReducedRegions = network.Regions{
		DigitalOcean: DOReducedRegions,
		Vultr:        VultrReducedRegions,
	}

	MinimalRegions = network.Regions{
		DigitalOcean: DOMinimalRegions,
		Vultr:        VultrMinimalRegions,
	}

	TestRegions = network.Regions{
		DigitalOcean: DOTestRegions,
		// Vultr:        VultrTestRegions,
		Linode: LinodeTestRegions,
	}
)
