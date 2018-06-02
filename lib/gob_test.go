package lib

import (
    "fmt"
    "testing"
)

type Post struct {
    Id   int
    Name string
    Job  string
}

func Test_gob_store_load(t *testing.T) {
    fmt.Println("gob文件读写")
    post := Post{
        Id:   1,
        Name: "heige",
        Job:  "goer",
    }
    StoreGobData(post, "post_gob.md")

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
