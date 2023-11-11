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
	AnnualEarningsBreakdown BreakDownJSON `json:"annualEarningsBreakdown" db:"annual_earnings_breakdown"`
	MonthlyEarningsBreakdown BreakDownJSON `json:"monthlyEarningsBreakdown" db:"monthly_earnings_breakdown"`
	WeeklyEarningsBreakdown BreakDownJSON `json:"weeklyEarningsBreakdown" db:"weekly_earnings_breakdown"`
	MonthlyWalkIns []WalkInData `json:"monthlyWalkIns" db:"monthly_walk_ins"`
	WeeklyWalkIns []WalkInData`json:"weeklyWalkIns" db:"weekly_walk_ins"`
}

type ClientDashboardData struct {
	Reservations int `json:"reservations" db:"reservations"`
	Packages int `json:"packages" db:"packages"`
	MembershipRequests int `json:"membershipRequests" db:"membership_requests"`
	CoachAppointments int `json:"coachAppointments" db:"coach_appointments"`
	AnnualExpendituresBreakdown BreakDownJSON `json:"annualExpendituresBreakdown" db:"annual_expenditures_breakdown"`
	MonthlyExpendituresBreakdown BreakDownJSON `json:"monthlyExpendituresBreakdown" db:"monthly_expenditures_breakdown"`
	WeeklyExpendituresBreakdown BreakDownJSON `json:"weeklyExpendituresBreakdown" db:"weekly_expenditures_breakdown"`
	WalkIns []ClientLog `json:"walkIns"`
}
type WalkInData struct {
	Total  int `json:"total" db:"total"`
	Date string `json:"date" db:"date"`
}

type BreakDown struct {
	WalkIn float64 `json:"walkIn" db:"walk_in"`
	Package float64 `json:"package" db:"package"`
	Membership float64 `json:"membership" db:"membership"`
}

type BreakDownJSON struct {
	BreakDown
}

func (instance *BreakDownJSON) Scan(value interface{}) error {
	val, valid := value.([]byte)
	if valid {
		unmarshalErr := json.Unmarshal(val, instance)
		if unmarshalErr != nil {
			*instance = BreakDownJSON{}
		}
	} else {
		*instance = BreakDownJSON{}
	}
	return nil

}
func (copy BreakDownJSON) Value(value interface{}) (driver.Value, error) {
	return copy, nil
}