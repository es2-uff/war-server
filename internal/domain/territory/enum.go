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

var TerritoryAdjacencyMap = map[TerritoryID][]TerritoryID{
	// Africa
	Algeria:     {Egypt, Sudan, Congo, Portugal, Brazil},
	Egypt:       {Algeria, Sudan, MiddleEast, Poland, Portugal},
	Sudan:       {Algeria, Egypt, Congo, SouthAfrica, Madagascar},
	Congo:       {Algeria, Sudan, SouthAfrica},
	SouthAfrica: {Sudan, Congo, Madagascar},
	Madagascar:  {Sudan, SouthAfrica},

	// Europe
	England:  {Iceland, Sweden, Germany, Portugal},
	Iceland:  {England, Sweden, Greenland},
	Sweden:   {England, Iceland, Moscow, Germany, Poland},
	Moscow:   {Sweden, Poland, Aral, Omsk},
	Germany:  {England, Sweden, Poland, Portugal},
	Poland:   {Sweden, Moscow, Germany, MiddleEast},
	Portugal: {England, Germany, Algeria, Brazil},

	// Asia
	MiddleEast: {Egypt, Poland, India, Aral},
	India:      {MiddleEast, Vietnam, China, Aral},
	Vietnam:    {India, China, Borneo},
	China:      {India, Vietnam, Mongolia, Vladvostok},
	Aral:       {Moscow, MiddleEast, India, Omsk},
	Omsk:       {Moscow, Aral, Dudinka, Mongolia},
	Dudinka:    {Omsk, Siberia, Mongolia, Mackenzie},
	Siberia:    {Dudinka, Tchita, Vladvostok, Alaska},
	Tchita:     {Siberia, Mongolia, Vladvostok},
	Mongolia:   {China, Omsk, Dudinka, Tchita, Japan},
	Japan:      {Mongolia, Vladvostok},
	Vladvostok: {China, Siberia, Tchita, Japan},

	// Oceania
	Australia: {NewGuinea, Sumatra, Borneo},
	NewGuinea: {Australia, Sumatra, Borneo},
	Sumatra:   {Australia, NewGuinea, Borneo},
	Borneo:    {Australia, NewGuinea, Sumatra, Vietnam},

	// South America
	Brazil:    {Algeria, Argentina, Chile, Colombia},
	Argentina: {Brazil, Chile},
	Chile:     {Brazil, Argentina, Colombia},
	Colombia:  {Brazil, Chile, Mexico},

	// North America
	Mexico:     {Colombia, California, NewYork},
	California: {Mexico, NewYork, Ottawa, Vancouver},
	NewYork:    {Mexico, California, Labrador, Ottawa},
	Labrador:   {NewYork, Ottawa, Greenland},
	Ottawa:     {California, NewYork, Labrador, Vancouver, Mackenzie},
	Vancouver:  {California, Ottawa, Mackenzie, Alaska},
	Mackenzie:  {Ottawa, Vancouver, Alaska, Dudinka, Greenland},
	Alaska:     {Vancouver, Mackenzie, Siberia},
	Greenland:  {Iceland, Labrador, Mackenzie},
}
