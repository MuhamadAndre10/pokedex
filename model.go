package main

type EncounterResponse struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	BaseExperience int        `json:"base_experience"`
	Height         int        `json:"height"`
	Weight         int        `json:"weight"`
	Stats          []StatInfo `json:"stats"`
	Types          []Types    `json:"types"`
}

type Types struct {
	Type Type `json:"type"`
}

type Type struct {
	Name string `json:"name"`
}

type StatInfo struct {
	BaseStat int  `json:"base_stat"`
	Stat     Stat `json:"stat"`
}

type Stat struct {
	Name string `json:"name"`
}
