package models

type CharacterData struct {
	General  map[string]string `json:"general"`
	Physical map[string]string `json:"physical"`
	Series   map[string]string `json:"series"`
}

type Character struct {
	URL       string        `json:"url"`
	Data      CharacterData `json:"data"`
	Abilities Abilities     `json:"abilities"`
}

// Characters holds the list of characters.
type Characters struct {
	Characters []Character `json:"characters"`
}

func NewCharacterData() CharacterData {
	return CharacterData{
		General:  make(map[string]string),
		Physical: make(map[string]string),
		Series:   make(map[string]string),
	}
}

func NewCharacter(url string) *Character {
	return &Character{
		URL:       url,
		Data:      NewCharacterData(),
		Abilities: make(map[string]string),
	}
}

func (c *Character) AddGeneralData(key string, value string) {
	c.Data.General[key] = value
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
