package psp_elastic_document

type PspNestedPokemon struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	PokedexNumber int32  `json:"pokedex_number"`
	Type1         string `json:"type_1"`
	Type2         string `json:"type_2"`
}

type PspNestedStock struct {
	Id     int64  `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

type PspElasticDocument struct {
	Id           int64            `json:"id"`
	Pokemon      PspNestedPokemon `json:"pokemon"`
	Stock        PspNestedStock   `json:"stock"`
	ActiveSeason bool             `json:"active_season"`
}
