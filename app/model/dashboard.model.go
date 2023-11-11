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
	MonthlyEarningsBreakdown EarningsBreakDownJSON `json:"monthlyEarningsBreakdown" db:"monthly_earnings_breakdown"`
	WeeklyEarningsBreakdown EarningsBreakDownJSON `json:"weeklyEarningsBreakdown" db:"weekly_earnings_breakdown"`
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