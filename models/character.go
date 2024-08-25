package models

type Character struct {
	URL       string            `json:"url"`
	Data      map[string]string `json:"data"`
	Abilities Abilities         `json:"abilities"`
}

// Characters holds the list of characters.
type Characters struct {
	Characters []Character `json:"characters"`
}

func NewCharacter(url string) *Character {
	return &Character{
		URL:       url,
		Data:      make(map[string]string),
		Abilities: make(map[string]string),
	}
}

func (c *Character) AddData(key string, value string) {
	c.Data[key] = value
}

func (c *Character) AddAbility(key string, value string) {
	c.Abilities[key] = value
}

func (c *Character) AddAbilities(abilities Abilities) {
	for key, value := range abilities {
		c.Abilities[key] = value
	}
}
