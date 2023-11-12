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