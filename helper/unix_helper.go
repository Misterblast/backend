package helper

import (
	"strconv"
	"time"
)

func FormatUnixTime(timeUnix string) string {
	timeUnixInt, _ := strconv.ParseInt(timeUnix, 10, 64)
	return time.Unix(timeUnixInt, 0).Format("2006-01-02 15:04:05")
}
