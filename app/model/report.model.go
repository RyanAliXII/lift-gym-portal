package model

import "time"

type ReportConfig struct {
	DateRange []time.Time `json:"dateRange"`
}