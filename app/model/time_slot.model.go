package model

import (
	"fmt"
	"lift-fitness-gym/app/db"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type TimeSlot  struct {
	Id int `json:"id"`
	StartTime string `json:"startTime" db:"start_time"`
	EndTime string `json:"endTime" db:"end_time"`
	MaxCapacity int `json:"maxCapacity" db:"max_capacity"`
	Available int `json:"available" db:"available"`
	Booked int `json:"booked,omitempty" db:"booked"`
	Model 
}
func (m TimeSlot)RemoveSecondsInTime() (string, string,error) {
	layout := "15:00:00"

    t, err := time.Parse(layout, m.StartTime)
    if err != nil {
        return "", "", err
    }
    startTime := t.Format("15:04")

	t, err = time.Parse(layout, m.EndTime  )
    if err != nil {
        return "", "", err
    }
	endTime := t.Format("15:04")   

	return startTime, endTime, nil

}
func validateTime(Value interface{}) error {
	layout := "15:04"
	timeStr, ok := Value.(string)
	if !ok {
		return fmt.Errorf("invalid time")
	}
	_, err := time.Parse(layout, timeStr)
	if err != nil {
		return fmt.Errorf("invalid time")
	}
	return nil
}

func ValidateEndTime (startTimeStr string) validation.RuleFunc {
	return func (value interface{}) error  {
		layout := "15:04"
		endTimeStr, ok := value.(string)
		if !ok {
			return fmt.Errorf("invalid time")
		}
		endTime, err := time.Parse(layout, endTimeStr)
		if err != nil {
			return fmt.Errorf("invalid time")
		}
		startTime, err := time.Parse(layout, startTimeStr)
		if err != nil {
			return fmt.Errorf("invalid time")
		}

		if(endTime.Before(startTime) || endTime.Equal(startTime)){
			return fmt.Errorf("end time must be after start time")
		}
		return nil
	}
}
func validateStartTimeIfTaken(value interface{}) error{
	startTime, _ := value.(string)
	db := db.GetConnection()
	recordCount := 1
	err := db.Get(&recordCount,"SELECT COUNT(1) as recordCount from time_slot where deleted_at is null and start_time = ?", startTime)
	if err != nil {
		return fmt.Errorf("start time already exists")
	}
	if recordCount >= 1 {
		return fmt.Errorf("start time already exists")
	}
	return nil
}
func validateEndTimeIfTaken(value interface{}) error{
	endTime, _ := value.(string)
	db := db.GetConnection()
	recordCount := 1
	err := db.Get(&recordCount,"SELECT COUNT(1) as recordCount from time_slot where deleted_at is null and end_time = ?", endTime)
	if err != nil {
		return fmt.Errorf("start time already exists")
	}
	if recordCount >= 1 {
		return fmt.Errorf("start time already exists")
	}
	return nil
}
func validateStartTimeIfTakenOnUpdate (id int) validation.RuleFunc{
	return func (value interface{}) error {
		startTime, _ := value.(string)
		db := db.GetConnection()
		recordCount := 1
		err := db.Get(&recordCount,"SELECT COUNT(1) as recordCount from time_slot where deleted_at is null and start_time = ? and id != ?", startTime, id)
		if err != nil {
			return fmt.Errorf("start time already exists")
		}
		if recordCount >= 1 {
			return fmt.Errorf("start time already exists")
		}
		return nil
	}
}
func validateEndTimeIfTakenOnUpdate (id int) validation.RuleFunc{
	return func (value interface{}) error {
		endTime, _ := value.(string)
		db := db.GetConnection()
		recordCount := 1
		err := db.Get(&recordCount,"SELECT COUNT(1) as recordCount from time_slot where deleted_at is null and end_time = ? and id != ?", endTime, id)
		if err != nil {
			return fmt.Errorf("start time already exists")
		}
		if recordCount >= 1 {
			return fmt.Errorf("start time already exists")
		}
			return nil
	}
}

func validateTimeSlotOverlap (startTime string, endTime string) error {
	db := db.GetConnection()
	recordCount := 1
	err := db.Get(&recordCount,"SELECT COUNT(1) FROM time_slot where deleted_at is null and start_time >= ? and end_time <= ?", startTime, endTime)
	if err != nil {
		return fmt.Errorf("time slot overlaps with current defined slots")
	}
	if recordCount >= 1 {
		return fmt.Errorf("time slot overlaps with current defined slots")
	}
	return nil
} 
func validateTimeSlotOverlapOnUpdate(startTime string, endTime string, id int) error {
	db := db.GetConnection()
	recordCount := 1
	err := db.Get(&recordCount,"SELECT COUNT(1) FROM time_slot where deleted_at is null and start_time >= ? and end_time <= ? and id != ?", startTime, endTime, id)
	if err != nil {
		return fmt.Errorf("time slot overlaps with current defined slots")
	}
	if recordCount >= 1 {
		return fmt.Errorf("time slot overlaps with current defined slots")
	}
	return nil
} 
func (m  TimeSlot) Validate() (error, map[string]string) {
	err := validateTimeSlotOverlap(m.StartTime, m.EndTime)
	if err != nil {
		return err, map[string]string{
			"startTime": err.Error(),
			"endTime": err.Error(),
		}
	}
	return m.Model.ValidationRules(&m, 
		validation.Field(&m.StartTime, validation.Required.Error("Start time is required."), validation.By(validateTime), validation.By(validateStartTimeIfTaken)), 
		validation.Field(&m.EndTime, validation.Required.Error("End time is required."), 
		validation.By(validateTime), validation.By(ValidateEndTime(m.StartTime)), validation.By(validateEndTimeIfTaken)),
		validation.Field(&m.MaxCapacity, validation.Required.Error("Max capacity is required."), validation.Min(1).Error("Max capacity must be atleast 1.")),
	)
}


func (m  TimeSlot) ValidateOnUpdate() (error, map[string]string) {
	err := validateTimeSlotOverlapOnUpdate(m.StartTime, m.EndTime, m.Id)
	if err != nil {
		return err, map[string]string{
			"startTime": err.Error(),
			"endTime": err.Error(),
		}
	}
	return m.Model.ValidationRules(&m, 
		validation.Field(&m.StartTime, validation.Required.Error("Start time is required."), validation.By(validateTime), validation.By(validateStartTimeIfTakenOnUpdate(m.Id))), 
		validation.Field(&m.EndTime, validation.Required.Error("End time is required."), 
		validation.By(validateTime), validation.By(ValidateEndTime(m.StartTime)), validation.By(validateEndTimeIfTakenOnUpdate(m.Id))),
		validation.Field(&m.MaxCapacity, validation.Required.Error("Max capacity is required."), validation.Min(1).Error("Max capacity must be atleast 1.")),
	)
}

type TimeSelection struct {
	Label string `json:"label"`
	Value string `json:"value"`
}
type timeSelections []TimeSelection

var Selections timeSelections  = timeSelections{
	{ Label: "12:00 AM", Value: "00:00" },
	{ Label: "1:00 AM", Value: "01:00" },
	{ Label: "2:00 AM", Value: "02:00" },
	{ Label: "3:00 AM", Value: "03:00" },
	{ Label: "4:00 AM", Value: "04:00" },
	{ Label: "5:00 AM", Value: "05:00" },
	{ Label: "6:00 AM", Value: "06:00" },
	{ Label: "7:00 AM", Value: "07:00" },
	{ Label: "8:00 AM", Value: "08:00" },
	{ Label: "9:00 AM", Value: "09:00" },
	{ Label: "10:00 AM", Value: "10:00" },
	{ Label: "11:00 AM", Value: "11:00" },
	{ Label: "12:00 PM", Value: "12:00" },
	{ Label: "1:00 PM", Value: "13:00" },
	{ Label: "2:00 PM", Value: "14:00" },
	{ Label: "3:00 PM", Value: "15:00" },
	{ Label: "4:00 PM", Value: "16:00" },
	{ Label: "5:00 PM", Value: "17:00" },
	{ Label: "6:00 PM", Value: "18:00" },
	{ Label: "7:00 PM", Value: "19:00" },
	{ Label: "8:00 PM", Value: "20:00" },
	{ Label: "9:00 PM", Value: "21:00" },
	{ Label: "10:00 PM", Value: "22:00" },
	{ Label: "11:00 PM", Value: "23:00" },
}
func NewTimeSelection() timeSelections{
	return  Selections
}
// this functions assumes that the time slots that have been passed is sorted by TimeSlot.StartTime in ascending order.
func (selections  timeSelections)RemoveSelectedSelections(timeSlots []TimeSlot) timeSelections{
	slotsMap := make(map[string]TimeSlot)
	newSelections := make([]TimeSelection, 0)
	//convert list of time slots to map and use start time as key.
	for _, slot := range timeSlots{
		start, end , err := slot.RemoveSecondsInTime() // convert start time and end time to 24:00:00 to 24:00.
		if err != nil {
			return Selections
		}
		slot.StartTime = start
		slot.EndTime = end
		slotsMap[start] =  slot
	}
	if(len(timeSlots) == 0 ){
		return Selections
	}
	endTime := ""
	/*
	check each time slot is defined.
	if time slot is not defined then add it to selection
	if time slot is defined, don't add the time slot from start time to end time.
	*/
	for i := 0; i < len(selections); i++{
		_, startExist := slotsMap[selections[i].Value]
		if startExist {
			endTime = slotsMap[selections[i].Value].EndTime
			continue
		}
		if(endTime ==  selections[i].Value){
			endTime = ""
		}
		if(endTime != ""){
			continue
		}
		newSelections = append(newSelections, selections[i])
	}
	return newSelections
}


