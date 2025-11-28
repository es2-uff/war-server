package territory

type Region int
type Shape int
type TerritoryID int

const (
	Circle Shape = iota
	Square
	Triangle
)

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

var TerritoryShapeMap = map[TerritoryID]Shape{
	Algeria:     Circle,
	Egypt:       Square,
	Sudan:       Triangle,
	Congo:       Circle,
	SouthAfrica: Square,
	Madagascar:  Triangle,

	England:  Circle,
	Iceland:  Square,
	Sweden:   Triangle,
	Moscow:   Circle,
	Germany:  Square,
	Poland:   Triangle,
	Portugal: Square,

	MiddleEast: Circle,
	India:      Square,
	Vietnam:    Triangle,
	China:      Circle,
	Aral:       Square,
	Omsk:       Triangle,
	Dudinka:    Circle,
	Siberia:    Square,
	Tchita:     Triangle,
	Mongolia:   Circle,
	Japan:      Square,
	Vladvostok: Triangle,

	Australia: Circle,
	NewGuinea: Square,
	Sumatra:   Triangle,
	Borneo:    Square,

	Brazil:    Circle,
	Argentina: Square,
	Chile:     Triangle,
	Colombia:  Square,

	Mexico:     Circle,
	California: Square,
	NewYork:    Triangle,
	Labrador:   Circle,
	Ottawa:     Square,
	Vancouver:  Triangle,
	Mackenzie:  Circle,
	Alaska:     Square,
	Greenland:  Triangle,
}

var TerritoryNameMap = map[TerritoryID]string{
	// Africa
	Algeria:     "Argélia",
	Egypt:       "Egito",
	Sudan:       "Sudão",
	Congo:       "Congo",
	SouthAfrica: "África do Sul",
	Madagascar:  "Madagascar",

	// Europe
	England:  "Inglaterra",
	Iceland:  "Islândia",
	Sweden:   "Suécia",
	Moscow:   "Moscou",
	Germany:  "Alemanha",
	Poland:   "Polônia",
	Portugal: "Portugal",

	// Asia
	MiddleEast: "Oriente Médio",
	India:      "Índia",
	Vietnam:    "Vietnã",
	China:      "China",
	Aral:       "Aral",
	Omsk:       "Omsk",
	Dudinka:    "Dudinka",
	Siberia:    "Sibéria",
	Tchita:     "Tchita",
	Mongolia:   "Mongólia",
	Japan:      "Japão",
	Vladvostok: "Vladivostok",

	// Oceania
	Australia: "Austrália",
	NewGuinea: "Nova Guiné",
	Sumatra:   "Sumatra",
	Borneo:    "Bornéu",

	// South America
	Brazil:    "Brasil",
	Argentina: "Argentina",
	Chile:     "Chile",
	Colombia:  "Colômbia",

	// North America
	Mexico:     "México",
	California: "Califórnia",
	NewYork:    "Nova York",
	Labrador:   "Labrador",
	Ottawa:     "Ottawa",
	Vancouver:  "Vancouver",
	Mackenzie:  "Mackenzie",
	Alaska:     "Alasca",
	Greenland:  "Groenlândia",
}
