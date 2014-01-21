package bjtime

import (
	"time"
)

var (
	cnTzAjust, _ = time.LoadLocation("Asia/Shanghai")
)

// timestamp of UTC to beijing time
func TsToString(ts int) string {
	t := time.Unix(int64(ts), 0)
	return t.In(cnTzAjust).Format("01-02 15:04:05")
}

func TimeToString(t time.Time) string {
	return t.In(cnTzAjust).Format("01-02 15:04:05")
}

func NowBj() time.Time {
	return time.Now().In(cnTzAjust)
}
