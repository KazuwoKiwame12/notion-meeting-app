package model

type Template struct {
	Parent     `json:"parent"`
	Properties `json:"properties"`
	Children   []Block `json:"children"`
}

type Parent struct {
	DatabaseID string `json:"database_id"`
}

type Properties struct {
	Name struct {
		Title []struct {
			Text `json:"text"`
		} `json:"title"`
	}
}
