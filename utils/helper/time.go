package helper

import (
	"time"

	"github.com/sirupsen/logrus"
)

var (
	LocationJakarta *time.Location
)

func init() {
	LocationJakarta, _ = time.LoadLocation("Asia/Jakarta")
}

func GetNowTime() time.Time {
	return time.Now().In(LocationJakarta)
}

func StringToDateTime(s string) time.Time {
	parse, err := time.ParseInLocation(time.DateTime, s, LocationJakarta)
	if err != nil {
		logrus.Error(err)
	}

	return parse
}

func DateTimeToString(d time.Time) string {
	return d.Format(time.DateTime)
}
