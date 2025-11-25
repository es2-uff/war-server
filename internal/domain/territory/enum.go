package territory

type Region int
type TerritoryID int

const (
	Europe Region = iota
	Asia
	Africa
	Oceania
	SouthAmerica
	NorthAmerica
)

const (
	// Africa territories
	Algeria TerritoryID = iota
	Egypt
	Sudan
	Congo
	SouthAfrica
	Madagascar

	// Europe territories
	England
	Iceland
	Sweden
	Moscow
	Germany
	Poland
	Portugal

	// Asia territories
	MiddleEast
	India
	Vietnam
	China
	Aral
	Omsk
	Dudinka
	Siberia
	Tchita
	Mongolia
	Japan
	Vladvostok

	// Oceania territories
	Australia
	NewGuinea
	Sumatra
	Borneo

	// South America territories
	Brazil
	Argentina
	Chile
	Colombia

	// North America territories
	Mexico
	California
	NewYork
	Labrador
	Ottawa
	Vancouver
	Mackenzie
	Alaska
	Greenland
)

var TerritoryRegionMap = map[TerritoryID]Region{
	// Africa
	Algeria:     Africa,
	Egypt:       Africa,
	Sudan:       Africa,
	Congo:       Africa,
	SouthAfrica: Africa,
	Madagascar:  Africa,

	// Europe
	England:  Europe,
	Iceland:  Europe,
	Sweden:   Europe,
	Moscow:   Europe,
	Germany:  Europe,
	Poland:   Europe,
	Portugal: Europe,

	// Asia
	MiddleEast: Asia,
	India:      Asia,
	Vietnam:    Asia,
	China:      Asia,
	Aral:       Asia,
	Omsk:       Asia,
	Dudinka:    Asia,
	Siberia:    Asia,
	Tchita:     Asia,
	Mongolia:   Asia,
	Japan:      Asia,
	Vladvostok: Asia,

	// Oceania
	Australia: Oceania,
	NewGuinea: Oceania,
	Sumatra:   Oceania,
	Borneo:    Oceania,

	// South America
	Brazil:    SouthAmerica,
	Argentina: SouthAmerica,
	Chile:     SouthAmerica,
	Colombia:  SouthAmerica,

	// North America
	Mexico:     NorthAmerica,
	California: NorthAmerica,
	NewYork:    NorthAmerica,
	Labrador:   NorthAmerica,
	Ottawa:     NorthAmerica,
	Vancouver:  NorthAmerica,
	Mackenzie:  NorthAmerica,
	Alaska:     NorthAmerica,
	Greenland:  NorthAmerica,
}

var AllTerritories = []TerritoryID{
	// Africa
	Algeria, Egypt, Sudan, Congo, SouthAfrica, Madagascar,
	// Europe
	England, Iceland, Sweden, Moscow, Germany, Poland, Portugal,
	// Asia
	MiddleEast, India, Vietnam, China, Aral, Omsk, Dudinka, Siberia, Tchita, Mongolia, Japan, Vladvostok,
	// Oceania
	Australia, NewGuinea, Sumatra, Borneo,
	// South America
	Brazil, Argentina, Chile, Colombia,
	// North America
	Mexico, California, NewYork, Labrador, Ottawa, Vancouver, Mackenzie, Alaska, Greenland,
}
