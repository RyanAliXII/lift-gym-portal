package model

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type TimeSlot  struct {
	Id int `json:"id"`
	StartTime string `json:"startTime" db:"start_time"`
	EndTime string `json:"endTime" db:"end_time"`
	MaxCapacity int `json:"maxCapacity" db:"max_capacity"`
	Model 
}

func validateTime(value interface{}) error {
	layout := "15:04"
	timeStr, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid time")
	}
	_, err := time.Parse(layout, timeStr)
	if err != nil {
		return fmt.Errorf("invalid time")
	}
	return nil
}

func (m  TimeSlot) Validate() (error, map[string]string) {
	return m.Model.ValidationRules(&m, 
		validation.Field(&m.StartTime, validation.Required.Error("Start time is required."), validation.By(validateTime)), 
		validation.Field(&m.EndTime, validation.Required.Error("End time is required."), validation.By(validateTime)),
		validation.Field(&m.MaxCapacity, validation.Required.Error("Max capacity is required."), validation.Min(1).Error("Max capacity must be atleast 1.")),
	)
}