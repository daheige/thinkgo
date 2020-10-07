package def

// EmptyStruct zero size, empty struct
type EmptyStruct struct{}

// EmptyObject to json空对象{}格式返回
type EmptyObject struct{}

// EmptyArray 兼容其他语言的[]空数组,一般在to json的时候转换为[]
type EmptyArray []struct{}

// H 对map[string]interface{}别名类型，简化书写
type H map[string]interface{}
