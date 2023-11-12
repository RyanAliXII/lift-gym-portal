package model

import (
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


type ReportData struct {
	Clients int `json:"clients" db:"clients"`
	Coaches int `json:"coaches" db:"coaches"`
	Members int `json:"members" db:"members"`
	InventoryItems int `json:"inventoryItems" db:"inventory_items"`
	Reservations int `json:"reservations" db:"reservations"`//
	MembershipRequests int `json:"membershipRequests" db:"membership_requests"`
	WalkIn []WalkInData `json:"walkIn" db:"walk_in"`
	PackageRequests int `json:"packageRequests" db:"package_requests"` 
	Earnings float64 `json:"earnings" db:"earnings"` 
	EarningsBreakdown BreakDownJSON `json:"earningsBreakdown" db:"earnings_breakdown"`
}


