package models

import (
	"time"
)

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role"`
}

type PickupPoint struct {
	ID               string       `json:"id"`
	RegistrationDate time.Time    `json:"registrationDate"`
	City             string       `json:"city"`
	Receptions       []*Reception `json:"receptions,omitempty"`
}

type Reception struct {
	ID            string     `json:"id"`
	DateTime      time.Time  `json:"dateTime"`
	Status        string     `json:"status"`
	PickupPointID string     `json:"pvzId"`
	Products      []*Product `json:"products,omitempty"`
}

type Product struct {
	ID          string    `json:"id"`
	DateTime    time.Time `json:"dateTime"`
	TypeProduct string    `json:"type"`
	ReceptionID string    `json:"receptionId"`
}