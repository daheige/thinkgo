package common

// EmptyStruct zero size, empty struct
type EmptyStruct struct{}

// EmptyArray 兼容其他语言的[]空数组,一般在tojson的时候转换为[]
type EmptyArray []struct{}

// H 对map[string]interface{}别名类型，简化书写
type H map[string]interface{}
