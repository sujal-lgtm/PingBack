package models

type Event struct{
	ID string `json:"id"`
	Source string `json:"source"`
	Payload string `json:"payload"`
}