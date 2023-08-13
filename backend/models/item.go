package models

type Item struct {
	Id       string `json:"Id"`
	Name     string `json:"Name"`
	Price    int    `json:"Price"`
	Stock    int    `json:"Stock"`
	Category int    `json:"Category"`
	Image    string `json:"Image"`
	Sold     int    `json:"Sold"`
}
