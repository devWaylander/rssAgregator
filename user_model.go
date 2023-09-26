package main

import "github.com/go-openapi/strfmt"

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	CreatedAt strfmt.DateTime `json:"created_at"`
	UpdatedAt strfmt.DateTime `json:"updated_at"`
}
