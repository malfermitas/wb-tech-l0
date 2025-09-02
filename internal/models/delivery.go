package models

type Delivery struct {
	Name    string `json:"name" fake:"{name}"`
	Phone   string `json:"phone" fake:"{phone}"`
	Zip     string `json:"zip" fake:"{zip}"`
	City    string `json:"city" fake:"{city}"`
	Address string `json:"address" fake:"{street}"`
	Region  string `json:"region" fake:"{state}"`
	Email   string `json:"email" fake:"{email}"`
}
