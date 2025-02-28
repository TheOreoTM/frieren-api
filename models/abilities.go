package models

type Abilities []Ability
type Ability struct {
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	SubAbilities []Ability `json:"subAbilities"`
}
