package models

import "time"

type StatusCheck struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status"`
}

type Response[T any] struct {
	Message    string    `json:"message"`
	StatusCode int       `json:"status"`
	Timestamp  time.Time `json:"timestamp"`
	// Device     *DeviceInfo `json:"device"`
	Data  *T        `json:"data,omitempty"`
	Error *APIError `json:"error,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Details string `json:"details,omitempty"`
}

type DeviceInfo struct {
	UserAgent string `json:"user_agent"`
	IPAddress string `json:"ip_address"`
}
