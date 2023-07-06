package config

import (
	"time"
)

var (
	Action          = "Action"
	Backend         = "Backend"
	RealIp          = "API-Real-IP"
)

var (
	MinimalTimeout  = 5 * time.Second
	Medium10Timeout = 10 * time.Second
	Medium15Timeout = 15 * time.Second
	LongestTimeout  = 20 * time.Second

	DefaultTimeout = LongestTimeout
)