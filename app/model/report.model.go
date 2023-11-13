package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type ReportConfig struct {
	DateRange [2]time.Time `json:"dateRange"`
}
func(config ReportConfig ) ToDateOnly()(start string, end string, parseErr error ){
	if(len(config.DateRange) < 2){
		return "","",fmt.Errorf("date range must be 2 : [0] for start [1] for end")
	}
	start = config.DateRange[0].Format(time.DateOnly)
	end = config.DateRange[1].Format(time.DateOnly)
	return start, end, nil
}


func (c ReportData) Value() (driver.Value, error) {
    b, err := json.Marshal(c)
    if err != nil {
        return nil, err
    }
    return string(b), nil
  
}
type ReportData struct {
	Id int64 `json:"id" db:"id"`
	Clients int `json:"clients" db:"clients"`
	StartDate string `json:"startDate" db:"start_date"`
	EndDate string `json:"endDate" db:"start_date"`
	Coaches int `json:"coaches" db:"coaches"`
	Members int `json:"members" db:"members"`
	InventoryItems int `json:"inventoryItems" db:"inventory_items"`
	Reservations int `json:"reservations" db:"reservations"`//
	MembershipRequests int `json:"membershipRequests" db:"membership_requests"`
	WalkIns []WalkInData `json:"walkIns" db:"walk_in"`
	PackageRequests int `json:"packageRequests" db:"package_requests"` 
	Earnings float64 `json:"earnings" db:"earnings"` 
	EarningsBreakdown BreakDownJSON `json:"earningsBreakdown" db:"earnings_breakdown"`
	PreparedBy string `json:"preparedBy" db:"prepared_by"`
}


type ReportDataJSON struct {
	ReportData
}

func (instance * ReportDataJSON) Scan(value interface{}) error {
	val, valid := value.([]byte)
	if valid {
		unmarshalErr := json.Unmarshal(val, instance)
		if unmarshalErr != nil {
			*instance = ReportDataJSON{}
		}
	} else {
		*instance = ReportDataJSON{}
	}
	return nil

}
func (copy ReportDataJSON) Value() (driver.Value, error) {
	return copy, nil
}