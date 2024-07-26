package main

type Regions struct {
	Vultr        map[string]int
	DigitalOcean map[string]int
}

var (
	FullRegions = Regions{
		DigitalOcean: DOFullRegions,
		Vultr:        VultrFullRegions,
	}

	HalfRegions = Regions{
		DigitalOcean: DOHalfRegions,
		Vultr:        VultrHalfRegions,
	}

	ReducedRegions = Regions{
		DigitalOcean: DOReducedRegions,
		Vultr:        VultrReducedRegions,
	}

	MinimalRegions = Regions{
		DigitalOcean: DOMinimalRegions,
		Vultr:        VultrMinimalRegions,
	}

	TestRegions = Regions{
		DigitalOcean: DOTestRegions,
		Vultr:        VultrTestRegions,
	}
)
