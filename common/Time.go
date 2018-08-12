package common

import (
	"time"
)

const tmFmtWithMS = "2006-01-02 15:04:05.999"
const tmFmtMissMS = "2006-01-02 15:04:05"
const tmFmtTime = "2006-01-02"

var TimeZone = "PRC" //默认时区设置

func SetTimeZone(zone string) {
	TimeZone = zone
}

// format a time.Time to string as 2006-01-02 15:04:05.999
func FormatTime(t time.Time) string {
	return t.Format(tmFmtWithMS)
}

// format a time.Time to string as 2006-01-02 15:04:05
// 将time.time转换为日期格式
func FormatTime19(t time.Time) string {
	return t.Format(tmFmtMissMS)
}

//获取时区loc
func GetLoc(zone string) *time.Location {
	loc, _ := time.LoadLocation(TimeZone)
	return loc
}

//当前本地时间
func GetCurrentLocalTime() string {
	return GetTimeByTimeZone(TimeZone)
}

func GetTimeByTimeZone(zone string) string {
	loc, _ := time.LoadLocation(zone)
	return time.Now().In(loc).Format(tmFmtMissMS)
}

// format time.Now() use FormatTime
func FormatNow() string {
	return FormatTime(time.Now())
}

// format time.Now().UTC() use FormatTime
func FormatUTC() string {
	return FormatTime(time.Now().UTC())
}

// parse a string to time.Time
func ParseTime(s string) (time.Time, error) {
	if len(s) == len(tmFmtMissMS) {
		return time.ParseInLocation(tmFmtMissMS, s, time.Local)
	}
	return time.ParseInLocation(tmFmtWithMS, s, time.Local)
}

// parse a string as "2006-01-02 15:04:05.999" to time.Time
func ParseTimeUTC(s string) (time.Time, error) {
	if len(s) == len(tmFmtMissMS) {
		return time.ParseInLocation(tmFmtMissMS, s, time.UTC)
	}
	return time.ParseInLocation(tmFmtWithMS, s, time.UTC)
}

// format a time.Time to number as 20060102150405999
func NumberTime(t time.Time) uint64 {
	y, m, d := t.Date()
	h, M, s := t.Clock()
	ms := t.Nanosecond() / 1000000
	return uint64(ms+s*1000+M*100000+h*10000000+d*1000000000) +
		uint64(m)*100000000000 + uint64(y)*10000000000000
}

// format time.Now() use NumberTime
func NumberNow() uint64 {
	return NumberTime(time.Now())
}

// format time.Now().UTC() use NumberTime
func NumberUTC() uint64 {
	return NumberTime(time.Now().UTC())
}

// parse a uint64 as 20060102150405999 to time.Time
func parseNumber(t uint64, tl *time.Location) (time.Time, error) {
	ns := int((t % 1000) * 1000000)
	t /= 1000
	s := int(t % 100)
	t /= 100
	M := int(t % 100)
	t /= 100
	h := int(t % 100)
	t /= 100
	d := int(t % 100)
	t /= 100
	m := time.Month(t % 100)
	y := int(t / 100)

	return time.Date(y, m, d, h, M, s, ns, tl), nil
}

func ParseNumber(t uint64) (time.Time, error) {
	return parseNumber(t, time.Local)
}

func ParseNumberUTC(t uint64) (time.Time, error) {
	return parseNumber(t, time.UTC)
}
