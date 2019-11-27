package gfile

import (
	"fmt"
	"log"
	"testing"
)

type Post struct {
	Id   int
	Name string
	Job  string
}

func TestGobStoreLoad(t *testing.T) {
	fmt.Println("gob文件读写")
	post := Post{
		Id:   1,
		Name: "heige313",
		Job:  "goer php",
	}
	err := StoreGobData(post, "post_gob.md")
	if err != nil {
		log.Println(err)
		return
	}

	//从文件中载入gob写入的文件内容到post
	var postData Post

	LoadGobData(&postData, "post_gob.md") //第一个参数接受postData的内存地址，因为loadData载入的数据会存入data中
	fmt.Println(postData)
	fmt.Println(postData.Id, postData.Name)

	//实现字符串的存取
	StoreGobData("fefefe", "test_gob.md")
	var str string
	LoadGobData(&str, "test_gob.md")
	fmt.Println(str)

	t.Log("success")
}

/**
$ go test -v
=== RUN   TestGobStoreLoad
gob文件读写
{1 heige313 goer php}
1 heige313
fefefe
--- PASS: TestGobStoreLoad (0.02s)
    gob_test.go:41: success
PASS
ok      github.com/daheige/thinkgo/gfile        0.022s
*/
