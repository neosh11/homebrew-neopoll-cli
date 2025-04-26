package config

import (
	"os"
)

// BaseURL is the API endpoint; override via NEOPOLL_BASE_URL env var if needed
var BaseURL = func() string {
	if v := os.Getenv("NEOPOLL_BASE_URL"); v != "" {
		return v
	}
	return "https://api-poll.lets.lol"
}()