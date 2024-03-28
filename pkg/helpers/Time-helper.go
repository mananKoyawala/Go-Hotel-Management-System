package helpers

import "time"

func GetTime() (time.Time, error) {
	return time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
}

func AddTime() (time.Time, error) {
	return time.Parse(time.RFC3339, time.Now().AddDate(0, 0, 1).Format(time.RFC3339))
}
