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

func (m  TimeSlot) Validate() (error, map[string]string) {
	return m.Model.ValidationRules(&m, 
		validation.Field(&m.StartTime, validation.Required.Error("Start time is required."), validation.By(validateTime)), 
		validation.Field(&m.EndTime, validation.Required.Error("End time is required."), validation.By(validateTime), validation.By(ValidateEndTime(m.StartTime))),
		validation.Field(&m.MaxCapacity, validation.Required.Error("Max capacity is required."), validation.Min(1).Error("Max capacity must be atleast 1.")),
	)
}


type TimeSelection struct {
	Label string `json:"label"`
	Value string `json:"value"`
}
type timeSelections []TimeSelection


func NewTimeSelection() timeSelections{
	return  Selections
}

func (selections  timeSelections)RemoveSelectedSelections(timeSlots []TimeSlot) timeSelections{
	slotsMap := make(map[string]string)
	newSelections := make([]TimeSelection, 0)
	for _, slot := range timeSlots{
		start, end , err := slot.RemoveSecondsInTime()
		
		if err != nil {
			return Selections
		}
		slotsMap[start] =  end
	}
	if(len(timeSlots) == 0 ){
		return Selections
	}
	endTime := ""
	for i := 0; i < len(selections); i++{
		_, startExist := slotsMap[selections[i].Value]
		if startExist {
			endTime = slotsMap[selections[i].Value]
			continue
		}
		if(endTime ==  selections[i].Value){
			endTime = ""
			continue
		}
		if(endTime != ""){
			continue
		}
		newSelections = append(newSelections, selections[i])
	}
	return newSelections
}
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
func getTimeIndexPosition () map[string]int {
	positions := make(map[string]int)
	for idx, selection := range Selections {
		positions[selection.Value] = idx
	}
	return positions
}

var timePositions = getTimeIndexPosition()



