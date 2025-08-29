package main

import "time"

type ServiceStatus struct {
	Service      string    `json:"service"`
	URL          string    `json:"url"`
	HTTPUp       bool      `json:"http_up"`
	HTTPCode     int       `json:"http_code"`
	ResponseTime int64     `json:"response_time"`
	DNSValid     bool      `json:"dns_valid"`
	DNSError     string    `json:"dns_error"`
	IsBlocked    bool      `json:"is_blocked"`
	LastChecked  time.Time `json:"last_checked"`
}

type Stats struct {
	TotalServices   int    `json:"total_services"`
	DownServices    int    `json:"down_services"`
	BlockedServices int    `json:"blocked_services"`
	LastCheck       string `json:"last_check"`
}
