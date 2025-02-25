package models

type Abilities []Ability
type Ability struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	SubItems    []Ability `json:"subItems"`
}
