package objective

import "es2.uff/war-server/internal/domain/territory"

type ObjectiveID int

const (
	// Conquistar na totalidade a EUROPA, a OCEANIA e mais um terceiro.
	ConquerEuropeOceaniaAndOne ObjectiveID = iota

	// Conquistar na totalidade a ÁSIA e a AMÉRICA DO SUL.
	ConquerAsiaSouthAmerica

	// Conquistar na totalidade a EUROPA, a AMÉRICA DO SUL e mais um terceiro.
	ConquerEuropeSouthAmericaAndOne

	// Conquistar 18 TERRITÓRIOS e ocupar cada um deles com pelo menos dois exércitos.
	Conquer18TerritoriesWith2Armies

	// Conquistar na totalidade a ÁSIA e a ÁFRICA.
	ConquerAsiaAfrica

	// Conquistar na totalidade a AMÉRICA DO NORTE e a ÁFRICA.
	ConquerNorthAmericaAfrica

	// Conquistar 24 TERRITÓRIOS à sua escolha.
	Conquer24Territories

	// Conquistar na totalidade a AMÉRICA DO NORTE e a OCEANIA.
	ConquerNorthAmericaOceania
)

type ObjectiveType int

const (
	RegionConquest ObjectiveType = iota
	TerritoryCount
)

type Objective struct {
	ID          ObjectiveID
	Type        ObjectiveType
	Description string
	// For RegionConquest objectives
	RequiredRegions []territory.Region
	// For objectives that require "and one more region"
	RequiresAdditionalRegion bool
	// For TerritoryCount objectives
	RequiredTerritoryCount int
	MinArmiesPerTerritory  int
}

var AllObjectives = []ObjectiveID{
	ConquerEuropeOceaniaAndOne,
	ConquerAsiaSouthAmerica,
	ConquerEuropeSouthAmericaAndOne,
	Conquer18TerritoriesWith2Armies,
	ConquerAsiaAfrica,
	ConquerNorthAmericaAfrica,
	Conquer24Territories,
	ConquerNorthAmericaOceania,
}

var ObjectiveDetails = map[ObjectiveID]Objective{
	ConquerEuropeOceaniaAndOne: {
		ID:                       ConquerEuropeOceaniaAndOne,
		Type:                     RegionConquest,
		Description:              "Conquistar na totalidade a EUROPA, a OCEANIA e mais um terceiro.",
		RequiredRegions:          []territory.Region{territory.Europe, territory.Oceania},
		RequiresAdditionalRegion: true,
	},
	ConquerAsiaSouthAmerica: {
		ID:              ConquerAsiaSouthAmerica,
		Type:            RegionConquest,
		Description:     "Conquistar na totalidade a ÁSIA e a AMÉRICA DO SUL.",
		RequiredRegions: []territory.Region{territory.Asia, territory.SouthAmerica},
	},
	ConquerEuropeSouthAmericaAndOne: {
		ID:                       ConquerEuropeSouthAmericaAndOne,
		Type:                     RegionConquest,
		Description:              "Conquistar na totalidade a EUROPA, a AMÉRICA DO SUL e mais um terceiro.",
		RequiredRegions:          []territory.Region{territory.Europe, territory.SouthAmerica},
		RequiresAdditionalRegion: true,
	},
	Conquer18TerritoriesWith2Armies: {
		ID:                     Conquer18TerritoriesWith2Armies,
		Type:                   TerritoryCount,
		Description:            "Conquistar 18 TERRITÓRIOS e ocupar cada um deles com pelo menos dois exércitos.",
		RequiredTerritoryCount: 18,
		MinArmiesPerTerritory:  2,
	},
	ConquerAsiaAfrica: {
		ID:              ConquerAsiaAfrica,
		Type:            RegionConquest,
		Description:     "Conquistar na totalidade a ÁSIA e a ÁFRICA.",
		RequiredRegions: []territory.Region{territory.Asia, territory.Africa},
	},
	ConquerNorthAmericaAfrica: {
		ID:              ConquerNorthAmericaAfrica,
		Type:            RegionConquest,
		Description:     "Conquistar na totalidade a AMÉRICA DO NORTE e a ÁFRICA.",
		RequiredRegions: []territory.Region{territory.NorthAmerica, territory.Africa},
	},
	Conquer24Territories: {
		ID:                     Conquer24Territories,
		Type:                   TerritoryCount,
		Description:            "Conquistar 24 TERRITÓRIOS à sua escolha.",
		RequiredTerritoryCount: 24,
		MinArmiesPerTerritory:  1,
	},
	ConquerNorthAmericaOceania: {
		ID:              ConquerNorthAmericaOceania,
		Type:            RegionConquest,
		Description:     "Conquistar na totalidade a AMÉRICA DO NORTE e a OCEANIA.",
		RequiredRegions: []territory.Region{territory.NorthAmerica, territory.Oceania},
	},
}
