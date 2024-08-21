package model

type PokedexEntriesResponse struct {
	Description []struct {
		Description string `json:"description"`
		Language    struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"language"`
	} `json:"descriptions"`
	Id           int    `json:"id"`
	IsMainSeries bool   `json:"is_main_series"`
	Name         string `json:"name"`
	Names        []struct {
		Language struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEntries []struct {
		EntryNumber    int `json:"entry_number"`
		PokemonSpecies struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"pokemon_species"`
	} `json:"pokemon_entries"`
}
