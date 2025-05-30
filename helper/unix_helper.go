package helper

import (
	"strconv"
	"time"
)

func FormatUnixTime(timeUnix string) string {
	timeUnixInt, err := strconv.ParseInt(timeUnix, 10, 64)
	if err != nil {
		return ""
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return ""
	}

	t := time.Unix(timeUnixInt, 0).In(loc)
	return t.Format("2006-01-02 15:04:05")
}
