package model

type TimeSlot  struct {
	Id int `json:"id"`
	StartTime string `json:"startTime" db:"start_time"`
	EndTime string `json:"endTime" db:"end_time"`
	MaxCapacity int `json:"maxCapacity" db:"max_capacity"`
}