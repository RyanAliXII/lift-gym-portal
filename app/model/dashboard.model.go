package model

import (
	"database/sql/driver"
	"encoding/json"
)

type AdminDashboardData struct {
	Clients int `json:"clients" db:"clients"`
	Members int `json:"members" db:"members"`
	MonthlyEarnings float64 `json:"monthlyEarnings" db:"monthly_earnings"`
	WeeklyEarnings float64 `json:"weeklyEarnings" db:"weekly_earnings"`
	AnnualEarnings float64 `json:"annualEarnings" db:"annual_earnings"`
	AnnualEarningsBreakdown EarningsBreakDownJSON `json:"annualEarningsBreakdown" db:"annual_earnings_breakdown"`
	MonthlyEarningsBreakdown EarningsBreakDownJSON `json:"monthlyEarningsBreakdown" db:"monthly_earnings_breakdown"`
	WeeklyEarningsBreakdown EarningsBreakDownJSON `json:"weeklyEarningsBreakdown" db:"weekly_earnings_breakdown"`
	MonthlyWalkIns []WalkInData `json:"monthlyWalkIns" db:"monthly_walk_ins"`
	WeeklyWalkIns []WalkInData`json:"weeklyWalkIns" db:"weekly_walk_ins"`
}


type ClientDashboardData struct {
	Reservations int `json:"reservations" db:"reservations"`
	Packages int `json:"packages" db:"packages"`
	CoachAppointments int `json:"coachAppointments" db:"coach_appointments"`
	AnnualEarningsBreakdown EarningsBreakDownJSON `json:"annualEarningsBreakdown" db:"annual_earnings_breakdown"`
	MonthlyEarningsBreakdown EarningsBreakDownJSON `json:"monthlyEarningsBreakdown" db:"monthly_earnings_breakdown"`
	WeeklyEarningsBreakdown EarningsBreakDownJSON `json:"weeklyEarningsBreakdown" db:"weekly_earnings_breakdown"`
	
}
type WalkInData struct {
	Total  int `json:"total" db:"total"`
	Date string `json:"date" db:"date"`
}

type EarningsBreakDown struct {
	WalkIn float64 `json:"walkIn" db:"walk_in"`
	Package float64 `json:"package" db:"package"`
	Membership float64 `json:"membership" db:"membership"`
}

type EarningsBreakDownJSON struct {
	EarningsBreakDown
}

func (instance *EarningsBreakDownJSON) Scan(value interface{}) error {
	val, valid := value.([]byte)
	if valid {
		unmarshalErr := json.Unmarshal(val, instance)
		if unmarshalErr != nil {
			*instance = EarningsBreakDownJSON{}
		}
	} else {
		*instance = EarningsBreakDownJSON{}
	}
	return nil

}
func (copy EarningsBreakDownJSON) Value(value interface{}) (driver.Value, error) {
	return copy, nil
}