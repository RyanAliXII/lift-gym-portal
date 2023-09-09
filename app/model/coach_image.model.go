package model

type CoachImage struct {
	Id      int    `json:"id" db:"od"`
	Path    string `json:"path" db:"path"`
	CoachId int    `json:"coachId" db:"coach_id"`
}