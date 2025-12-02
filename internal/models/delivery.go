package models

type Delivery struct {
	Name    string `json:"name" fake:"{firstname} {lastname}" validate:"required,min=2,max=100"`
	Phone   string `json:"phone" fake:"{phone}" validate:"required,e164"`
	Zip     string `json:"zip" fake:"{zip}" validate:"required,min=5,max=10"`
	City    string `json:"city" fake:"{city}" validate:"required,min=2,max=50"`
	Address string `json:"address" fake:"{streetaddress}" validate:"required,min=5,max=200"`
	Region  string `json:"region" fake:"{state}" validate:"required,min=2,max=50"`
	Email   string `json:"email" fake:"{email}" validate:"required,email"`
}
