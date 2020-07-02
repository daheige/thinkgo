package strlist

//StringList helper

import (
	"reflect"
	"sort"

	"github.com/daheige/thinkgo/grecover"
)

// StringList: str list
type StringList []string

// Len len
func (sl StringList) Len() int {
	return len(sl)
}

// Less less
func (sl StringList) Less(i, j int) bool {
	if len(sl[i]) == len(sl[j]) {
		return sl[i] < sl[j]
	}
	return len(sl[i]) < len(sl[j])
}

// Swap swap
func (sl StringList) Swap(i, j int) {
	sl[i], sl[j] = sl[j], sl[i]
}

// UniqueAdd unique add a string
func (sl StringList) UniqueAdd(token string) StringList {
	for _, v := range sl {
		if v == token {
			return sl
		}
	}

	for i, v := range sl {
		if v == "" {
			sl[i] = token
			return sl
		}
	}
	return append(sl, token)
}

// Delete del string
func (sl StringList) Delete(token string) int {
	count := 0
	for i, v := range sl {
		if v == token {
			sl[i] = ""
			count++
		}
	}
	return count
}

// IsEmpty is empty
func (sl StringList) IsEmpty() bool {
	for _, v := range sl {
		if v != "" {
			return false
		}
	}
	return true
}

// Count return sl count
func (sl StringList) Count() int {
	count := 0
	for _, v := range sl {
		if v != "" {
			count++
		}
	}
	return count
}

// StringMapKeys string map sort
func StringMapKeys(m interface{}) (res []string) {
	defer grecover.CheckPanic()

	keys := reflect.ValueOf(m).MapKeys()
	for _, v := range keys {
		res = append(res, v.String())
	}

	sort.Sort(StringList(res))
	return
}
