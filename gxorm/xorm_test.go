package gxorm

import (
	"log"
	"testing"
	"time"

	"xorm.io/xorm"
)

/**
* sql:
* CREATE DATABASE IF NOT EXISTS test default charset utf8mb4;
* create table user (id int primary key auto_increment,name varchar(200),age tinyint) engine=innodb;
* 模拟数据插入
* mysql> insert into user (name) values("xiaoming");
   Query OK, 1 row affected (0.11 sec)

   mysql> insert into user (name) values("hello");
   Query OK, 1 row affected (0.04 sec)
*/

type myUser struct {
	Id   int    `xorm:"pk autoincr"` // 定义的字段属性，要用空格隔开
	Name string `xorm:"varchar(200)"`
	Age  int    `xorm:"tinyint(3)"`
}

func (myUser) TableName() string {
	return "user"
}

func TestGXORM(t *testing.T) {
	var e *xorm.Engine
	log.Println(e == nil)

	dbConf := &DbConf{
		DbBaseConf: DbBaseConf{
			Ip:        "127.0.0.1",
			Port:      3306,
			User:      "root",
			Password:  "root1234",
			Database:  "test",
			ParseTime: true,
		},

		MaxIdleConns: 10,
		MaxOpenConns: 100,
		ShowSql:      true,
	}

	db, err := dbConf.NewEngine() // 设置数据库连接对象，并非真正连接，只有在用的时候才会建立连接
	if db == nil || err != nil {
		log.Println("db error")
		return
	}

	defer db.Close()

	log.Println("====master db===")
	user := &myUser{}
	has, err := db.Where("id = ?", 1).Get(user)
	log.Println(has, err)
	log.Println("user info: ", user.Id, user.Name)

	// 测试读写分离
	rwConf := &EngineGroupConf{
		Master: DbBaseConf{
			Ip:        "127.0.0.1",
			Port:      3306,
			User:      "root",
			Password:  "root1234",
			Database:  "test",
			ParseTime: true,
		},
		Slaves: []DbBaseConf{
			DbBaseConf{
				Ip:        "127.0.0.1",
				Port:      3306,
				User:      "test1",
				Password:  "root1234",
				Database:  "test",
				ParseTime: true,
			},
			DbBaseConf{
				Ip:        "127.0.0.1",
				Port:      3306,
				User:      "test2",
				Password:  "root1234",
				Database:  "test",
				ParseTime: true,
			},
		},
		MaxIdleConns: 10,
		MaxOpenConns: 100,
		ShowSql:      true,
		MaxLifetime:  200 * time.Second,
	}

	eg, err := rwConf.NewEngineGroup()
	if err != nil {
		log.Println("set read db engine error: ", err.Error())
		return
	}

	userInfo := &myUser{}
	has, err = eg.Where("id = ?", 1).Get(userInfo)
	log.Println("get id = 1 of userInfo: ", has, err)

	log.Println("=======engine select=========")
	user2 := &myUser{}
	has, err = eg.Where("id = ?", 3).Get(user2)
	log.Println(has, err)
	log.Println(user2)

	// 采用读写分离实现数据插入
	user4 := &myUser{
		Name: "xiaoxiao",
		Age:  12,
	}

	// 插入单条数据，多条数据请用Insert(user3,user4,user5)
	affectedNum, err := eg.InsertOne(user4)
	log.Println("affected num: ", affectedNum)
	log.Println("insert id: ", user4.Id)
	log.Println("err: ", err)

	log.Println("get on slave to query")
	user5 := &myUser{}
	log.Println(eg.Slave().Where("id = ?", 4).Get(user5))
}

/**
$ go test -v
=== RUN   TestGXORM
2019/11/12 21:45:29 true
2019/11/12 21:45:29 ====master db===
[xorm] [info]  2019/11/12 21:45:29.633242 [SQL] SELECT `id`, `name`, `age` FROM `user` WHERE (id = ?) LIMIT 1 []interface {}{1} - took: 49.72224ms
2019/11/12 21:45:29 true <nil>
2019/11/12 21:45:29 user info:  1 xiaoxiao
[xorm] [info]  2019/11/12 21:45:29.665330 [SQL] SELECT `id`, `name`, `age` FROM `user` WHERE (id = ?) LIMIT 1 []interface {}{1} - took: 31.500405ms
2019/11/12 21:45:29 get id = 1 of userInfo:  true <nil>
2019/11/12 21:45:29 =======engine select=========
[xorm] [info]  2019/11/12 21:45:29.665787 [SQL] SELECT `id`, `name`, `age` FROM `user` WHERE (id = ?) LIMIT 1 []interface {}{3} - took: 297.281µs
2019/11/12 21:45:29 true <nil>
2019/11/12 21:45:29 &{3 xiaoxiao 12}
[xorm] [info]  2019/11/12 21:45:29.715391 [SQL] INSERT INTO `user` (`name`,`age`) VALUES (?, ?) []interface {}{"xiaoxiao", 12} - took: 49.453763ms
2019/11/12 21:45:29 affected num:  1
2019/11/12 21:45:29 insert id:  106
2019/11/12 21:45:29 err:  <nil>
2019/11/12 21:45:29 get on slave to query
[xorm] [info]  2019/11/12 21:45:29.743923 [SQL] SELECT `id`, `name`, `age` FROM `user` WHERE (id = ?) LIMIT 1 []interface {}{4} - took: 28.349699ms
2019/11/12 21:45:29 true <nil>
--- PASS: TestGxorm (0.16s)
PASS
ok      github.com/daheige/thinkgo/gxorm        0.164s
*/
