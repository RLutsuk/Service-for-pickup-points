package models

import "time"

// swagger
type InputUser struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password,omitempty" example:"password"`
	Role     string `json:"role" example:"employee"`
}
type OutputUser struct {
	ID    string `json:"id" example:"28b0a78e-dee5-4b9c-9a3f-61ab78b2f483"`
	Email string `json:"email" example:"user@example.com"`
	Role  string `json:"role" example:"employee"`
}

type AuthUser struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password,omitempty" example:"password"`
}

type CreationInPickupPoint struct {
	City string `json:"city" example:"Москва"`
}

type CreationOutPickupPoint struct {
	ID               string    `json:"id" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
	RegistrationDate time.Time `json:"registrationDate" example:"2025-04-14T21:09:46.237Z"`
	City             string    `json:"city" example:"Москва"`
}

type InReception struct {
	PickupPointID string `json:"pvzId" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
}

type OutReception struct {
	ID            string    `json:"id" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
	DateTime      time.Time `json:"dateTime" example:"2025-04-14T21:13:18.396Z"`
	Status        string    `json:"status" example:"in_progress"`
	PickupPointID string    `json:"pvzId" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
}

type OutReceptionClosed struct {
	ID            string    `json:"id" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
	DateTime      time.Time `json:"dateTime" example:"2025-04-14T21:13:18.396Z"`
	Status        string    `json:"status" example:"close"`
	PickupPointID string    `json:"pvzId" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
}

type InProduct struct {
	TypeProduct string `json:"type" example:"электроника"`
	ReceptionID string `json:"receptionId" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
}

type OutProduct struct {
	ID          string    `json:"id" example:"2e292031-a998-4a7a-ae4c-12908941858f"`
	DateTime    time.Time `json:"dateTime" example:"2025-04-14T21:17:17.911Z"`
	TypeProduct string    `json:"type" example:"электроника"`
	ReceptionID string    `json:"receptionId" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
}
