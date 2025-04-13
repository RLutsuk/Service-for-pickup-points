package models

import (
	"time"
)

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type PickupPoint struct {
	ID               string    `json:"id"`
	RegistrationDate time.Time `json:"registrationDate"`
	City             string    `json:"city"`
}

type Reception struct {
	ID            string    `json:"id,omitempty"`
	DateTime      time.Time `json:"dateTime,omitempty"`
	Status        string    `json:"status,omitempty"`
	PickupPointID string    `json:"pvzId,omitempty"`
}

type Product struct {
	ID          string    `json:"id"`
	DateTime    time.Time `json:"dateTime"`
	TypeProduct string    `json:"type"`
	ReceptionID string    `json:"receptionId"`
}
