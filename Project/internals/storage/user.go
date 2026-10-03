package storage

import "time"

type User struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Car       []Car     `json:"car"`
}

type Car struct {
	ID      string `json:"id_car"`
	Model   string `json:"model"`
	CarType string `json:"car_type"`
	UserId  string `json:"user_id"`
}

type UserDTO struct {
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Phone     string   `json:"phone"`
	CarDTO    []CarDTO `json:"car_dto"`
}

type CarDTO struct {
	CarDTO_ID string `json:"car_dto_id"`
	Model     string `json:"model"`
	Car_type  string `json:"car_type"`
}
