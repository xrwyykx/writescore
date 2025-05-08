package utils

import "time"

func MarshalTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func StringToTime(t string) (resultTime time.Time, err error) {
	resultTime, err = time.Parse("2006-01-02 15:04:05", t)
	return
}
