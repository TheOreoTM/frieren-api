package models

type Location struct {
	URL  string `json:"url"`
	Data map[string]string
}

type Locations struct {
	Central  []string          `json:"central"`
	Nothern  *NothernLocations `json:"nothern"`
	Southern []string          `json:"southern"`
}

type NothernLocations struct {
	NothernPlateau    []string `json:"nothern_plateau"`
	ImperialTerritory []string `json:"imperial_territory"`
	Ende              []string `json:"ende"`
}

func NewLocation(url string) *Location {
	return &Location{
		URL:  url,
		Data: make(map[string]string),
	}
}

func NewLocations() *Locations {
	return &Locations{
		Central:  []string{},
		Nothern:  NewNorthernLocations(),
		Southern: []string{},
	}
}

func NewNorthernLocations() *NothernLocations {
	return &NothernLocations{
		NothernPlateau:    []string{},
		ImperialTerritory: []string{},
		Ende:              []string{},
	}
}
