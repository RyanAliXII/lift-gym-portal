package model

import (
	"encoding/json"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type CoachSchedule struct{
	Id int `json:"id" db:"id"`
	Date string `json:"date" db:"date"`
	Time string `json:"time" db:"time"`
	CoachId int `json:"coachId" db:"coach_id"`
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
		HHMMSS24h := "15:04:05"
		_, err := time.Parse(layout, m.Time)
		if err == nil {
			return nil
		}
		_, err = time.Parse(HHMMSS24h, m.Time)
		if err != nil {
			return fmt.Errorf("invalid time")
		}
		return nil
	}),
	))
	
}
type CoachScheduleJSON struct {
	CoachSchedule
}

func (instance *CoachScheduleJSON) Scan(value interface{}) error {
	val, valid := value.([]byte)
	if valid {
		unmarshalErr := json.Unmarshal(val, instance)
		if unmarshalErr != nil {
			*instance = CoachScheduleJSON{}
		}
	} else {
		*instance = CoachScheduleJSON{}
	}
	return nil

}