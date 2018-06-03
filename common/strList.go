package common

import (
	"reflect"
	"sort"
)

//str list
type StringList []string

func (this StringList) Len() int {
	return len(this)
}

func (this StringList) Less(i, j int) bool {
	if len(this[i]) == len(this[j]) {
		return this[i] < this[j]
	}
	return len(this[i]) < len(this[j])
}

func (this StringList) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func (this StringList) UniqueAdd(token string) StringList {
	for _, v := range this {
		if v == token {
			return this
		}
	}

	for i, v := range this {
		if v == "" {
			this[i] = token
			return this
		}
	}
	return append(this, token)
}

func (this StringList) Delete(token string) int {
	count := 0
	for i, v := range this {
		if v == token {
			this[i] = ""
			count++
		}
	}
	return count
}

func (this StringList) IsEmpty() bool {
	for _, v := range this {
		if v != "" {
			return false
		}
	}
	return true
}

func (this StringList) Count() int {
	count := 0
	for _, v := range this {
		if v != "" {
			count++
		}
	}
	return count
}

func StringMapKeys(m interface{}) (res []string) {
	defer CheckPanic()

	keys := reflect.ValueOf(m).MapKeys()
	for _, v := range keys {
		res = append(res, v.String())
	}

	sort.Sort(StringList(res))
	return
}
