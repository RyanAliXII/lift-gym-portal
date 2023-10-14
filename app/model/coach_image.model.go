package model

import (
	"database/sql/driver"
	"encoding/json"
)

type CoachImage struct {
	Id      int    `json:"id" db:"id"`
	Path    string `json:"path" db:"path"`
	CoachId int    `json:"coachId" db:"coach_id"`
}


type CoachImages []string 

func (instance *CoachImages) Scan(value interface{}) error {
	images := make([]string, 0)
	val, valid := value.([]byte)
	if valid {
		unmarshalErr := json.Unmarshal(val, instance)
		if unmarshalErr != nil {
			*instance = images
		}
	} else {
		*instance = images
	}
	return nil

}
func (copy CoachImages) Value(value interface{}) (driver.Value, error) {
	return copy, nil
}




