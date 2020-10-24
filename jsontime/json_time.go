// package jsontime fix time.Time  json encode/decode bug.
package jsontime

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Time 自定义数据结构Time
type Time time.Time

// NullToEmptyStr null是否转换为空字符串形式
var NullToEmptyStr bool

const (
	// TmFormat time format layout.
	TmFormat = "2006-01-02 15:04:05"
)

// UnmarshalJSON json decode Time.
func (t *Time) UnmarshalJSON(data []byte) (err error) {
	str := string(data)
	if str == `null` || str == "" || str == `""` {
		return nil
	}

	now, err := time.ParseInLocation(`"`+TmFormat+`"`, str, time.Local)
	*t = Time(now)
	return
}

// MarshalJSON json encode Time.
func (t Time) MarshalJSON() ([]byte, error) {
	// log.Println(t.String())
	if t.String() == "0001-01-01 00:00:00" {
		if NullToEmptyStr {
			return []byte(`""`), nil
		}

		return []byte(`null`), nil
	}

	b := make([]byte, 0, len(TmFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, TmFormat)
	b = append(b, '"')
	return b, nil
}

func (t Time) String() string {
	return time.Time(t).Format(TmFormat)
}

// ====================fix go gorm table field null bug======
// 在 gorm 中只重写 MarshalJSON 是不够的
// 因为 ORM 在插入记录、读取记录时需要的相应执行 Value 和 Scan 方法
// 需要引入 database/sql/driver 包。为了方便使用
// 可以定义一个 BaseModel 来替代 gorm.Model
/**
//具体model定义time.Time 用jsonTime.Time替代
type BaseModel struct {
    // gorm.Model
    ID        uint        `gorm:"primary_key" json:"id"`
    CreatedAt jsonTime.Time  `json:"createdAt"`
    UpdatedAt jsonTime.Time  `json:"updatedAt"`
    DeletedAt *jsonTime.Time `sql:"index" json:"-"`
}
*/

// Value insert timestamp into mysql need this function.
func (t Time) Value() (driver.Value, error) {
	var zeroTime time.Time

	if time.Time(t).UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}

	return t, nil
}

// Scan valueOf time.Time
func (t *Time) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = Time(value)
		return nil
	}

	return fmt.Errorf("can not convert %v to timestamp", v)
}
