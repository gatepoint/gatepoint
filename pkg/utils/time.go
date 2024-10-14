package utils

import "time"

func NowUTC() *time.Time {
	t := time.Now().UTC()
	return &t
}
