package helpers

import "time"

func GetTime() (time.Time, error) {
	return time.Parse(time.RFC3339, time.Now().Local().Format(time.RFC3339))
}
