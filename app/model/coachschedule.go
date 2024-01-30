package model

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type CoachSchedule struct{
	Date string `json:"date"`
	Time string `json:"time"`
	Model 
}
func (m * CoachSchedule) Validate()(map[string]string, error) {
	return m.Model.Validate(m, validation.Field(&m.Date, validation.Required.Error("Date is required."), validation.By(func(value interface{}) error {
		_, err := time.Parse(time.DateOnly, m.Date)
		if err != nil{
			return fmt.Errorf("invalid date")
		}
		return nil
	})),
	validation.Field(&m.Time, validation.Required.Error("Time is required."),
	validation.By(func(value interface{}) error {
		layout := "15:04"
		_, err := time.Parse(layout, m.Time)
		if err != nil {
			return fmt.Errorf("invalid time")
		}
		return nil
	}),
	))
	
}