package models

type CharacterData struct {
	Names    Names             `json:"names"`
	General  map[string]string `json:"general"`
	Physical map[string]string `json:"physical"`
	Series   map[string]string `json:"series"`
}

type Names struct {
	English  string `json:"english"`
	Japanese string `json:"japanese"`
	Romaji   string `json:"romaji"`
}

type Character struct {
	URL       string        `json:"url"`
	Name      string        `json:"name"`
	Data      CharacterData `json:"data"`
	Abilities Abilities     `json:"abilities"`
}

// Characters holds the list of characters.
type Characters struct {
	Characters []Character `json:"characters"`
}

func NewCharacterData() CharacterData {
	return CharacterData{
		Names: Names{
			English:  "",
			Japanese: "",
			Romaji:   "",
		},
		General:  make(map[string]string),
		Physical: make(map[string]string),
		Series:   make(map[string]string),
	}
}

func NewCharacter(url string) *Character {
	return &Character{
		Name:      "",
		URL:       url,
		Data:      NewCharacterData(),
		Abilities: make(map[string]string),
	}
}

func (c *Character) SetName(name string) {
	c.Name = name
	c.Data.Names.English = name
}

func (c *Character) AddGeneralData(key string, value string) {
	c.Data.General[key] = value
}

func (c *Character) AddPhysicalData(key string, value string) {
	c.Data.Physical[key] = value
}

func (c *Character) AddSeriesData(key string, value string) {
	c.Data.Series[key] = value
}

func (c *Character) AddAbility(key string, value string) {
	c.Abilities[key] = value
}

func (c *Character) AddAbilities(abilities Abilities) {
	for key, value := range abilities {
		c.Abilities[key] = value
	}
}
