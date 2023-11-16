package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)
type ContactUs struct {
	Id int `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
	Message string `json:"message" db:"value"`
	CreatedAt string `json:"createdAt" db:"created_at"`
	Model
}
func (m ContactUs) Validate () (map[string]string, error){
	return m.Model.Validate(&m,
		 validation.Field(&m.Name, validation.Required.Error("Name is required")),
		 validation.Field(&m.Email, validation.Required.Error("Email is required"), is.Email.Error("Invalid email address")),
		 validation.Field(&m.Message, validation.Required.Error("Message is required"), validation.Length(1, 255)))		 
}