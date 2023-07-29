package model

type Unit struct {
	Owner   int    `json:"id"`
	Elem    string `json:"elem"`
	Name    string `json:"name"`
	Constel int    `json:"constel"`
}
