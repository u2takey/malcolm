package model

type Plugin struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Version     string   `json:"version"`
	Identifier  string   `json:"identifier"`
	Options     []Option `json:"options"`
}
