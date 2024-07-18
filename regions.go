package main

type Regions struct {
	AWS map[string]int
	DO  map[string]int
	GCP map[string]int
}

var (
	DOFullRegions = map[string]int{
		"nyc1": 4, "nyc3": 4, "tor1": 4, "sfo2": 6, "sfo3": 6, "ams3": 11, "sgp1": 13, "lon1": 7, "fra1": 5,
		"blr1": 13, "syd1": 13,
	}

	DOHalfRegions = map[string]int{
		"nyc1": 3, "nyc3": 3, "sfo2": 3, "sfo3": 3, "ams3": 6, "sgp1": 7, "lon1": 4, "fra1": 3, "tor1": 3,
		"blr1": 7, "syd1": 7,
	}

	DOReducedRegions = map[string]int{
		"nyc1": 1, "nyc3": 2, "sfo2": 2, "sfo3": 1, "ams3": 2, "sgp1": 3, "lon1": 2, "fra1": 2, "tor1": 2,
		"blr1": 2, "syd1": 2,
	}

	DOMinimalRegions = map[string]int{
		"nyc3": 1, "sfo3": 1, "ams3": 1, "sgp1": 1, "lon1": 1, "tor1": 1,
		"blr1": 1, "syd1": 1,
	}

	DOTestRegions = map[string]int{
		"sfo3": 1, "sgp1": 1,
	}
)

var (
	FullRegions = Regions{
		DO: DOFullRegions,
	}

	HalfRegions = Regions{
		DO: DOHalfRegions,
	}

	ReducedRegions = Regions{
		DO: DOReducedRegions,
	}

	MinimalRegions = Regions{
		DO: DOMinimalRegions,
	}

	TestRegions = Regions{
		DO: DOTestRegions,
	}
)
