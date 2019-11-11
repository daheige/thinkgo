package gxorm

import (
	"log"
	"testing"

	"github.com/go-xorm/xorm"
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
	Id   int    `xorm:"pk autoincr"` //定义的字段属性，要用空格隔开
	Name string `xorm:"varchar(200)"`
	Age  int    `xorm:"tinyint(3)"`
}

func (myUser) TableName() string {
	return "user"
}

func TestGxorm(t *testing.T) {
	var e *xorm.Engine
	log.Println(e == nil)

	dbConf := &DbConf{

		Ip:        "127.0.0.1",
		Port:      3306,
		User:      "root",
		Password:  "root",
		Database:  "test",
		ParseTime: true,

		MaxIdleConns: 10,
		MaxOpenConns: 100,
		SqlCmd:       true,
		ShowExecTime: false,
	}

	db, err := dbConf.NewEngine() //设置数据库连接对象，并非真正连接，只有在用的时候才会建立连接
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

	//测试读写分离
	readConf := &DbConf{
		Ip:        "127.0.0.1",
		Port:      3306,
		User:      "test1",
		Password:  "1234",
		Database:  "test",
		ParseTime: true,

		MaxIdleConns: 10,
		MaxOpenConns: 100,
		SqlCmd:       true,
		ShowExecTime: true,
	}

	readDb, err := readConf.NewEngine()
	if err != nil {
		log.Println("set read db engine error: ", err.Error())
		return
	}

	defer readDb.Close()

	log.Println("===========read db of one=======")
	userInfo := &myUser{}
	has, err = readDb.Where("id = ?", 1).Get(userInfo)
	log.Println("read one db,get id = 1 of userInfo: ", has, err)

	//设置第二个读的实例
	readConf2 := &DbConf{
		Ip:        "127.0.0.1",
		Port:      3306,
		User:      "test2",
		Password:  "1234",
		Database:  "test",
		ParseTime: true,

		MaxIdleConns: 10,
		MaxOpenConns: 100,
		SqlCmd:       true,
		ShowExecTime: true,
	}

	readDb2, err := readConf2.NewEngine()
	if err != nil {
		log.Println("set read db engine error: ", err.Error())
		return
	}

	defer readDb2.Close()

	log.Println("=========slave two db====")
	userInfo2 := &myUser{}
	has, err = readDb2.Where("id = ?", 2).Get(userInfo2)
	log.Println("read two db user: ", has, err)

	//设置读写分离的引擎句柄
	// engineGroup, err := NewEngineGroup(db, readDb)
	// engineGroup, err := NewEngineGroup(db, readDb2)
	engineGroup, err := NewEngineGroup(db, readDb, readDb2)
	if err != nil {
		log.Println("set db engineGroup error: ", err.Error())
		return
	}

	defer engineGroup.Close() //关闭读写分离的连接对象

	log.Println("=======engine select=========")
	user2 := &myUser{}
	has, err = engineGroup.Where("id = ?", 3).Get(user2)
	log.Println(has, err)
	log.Println(user2)

	//采用读写分离实现数据插入
	user4 := &myUser{
		Name: "xiaoxiao",
		Age:  12,
	}

	affectedNum, err := engineGroup.InsertOne(user4) //插入单条数据，多条数据请用Insert(user3,user4,user5)
	log.Println("affected num: ", affectedNum)
	log.Println("insert id: ", user4.Id)
	log.Println("err: ", err)

	log.Println("get on slave to query")
	user5 := &myUser{}
	log.Println(engineGroup.Slave().Where("id = ?", 4).Get(user5))
}

/**
$ go test -v
=== RUN   TestGxorm
2019/11/12 00:01:21 true
2019/11/12 00:01:21 ====master db===
[xorm] [info]  2019/11/12 00:01:21.831136 [SQL] SELECT `id`, `name`, `age` FROM `user` WHERE (id = ?) LIMIT 1 []interface {}{1}
2019/11/12 00:01:21 true <nil>
2019/11/12 00:01:21 user info:  1 xiaoxiao
2019/11/12 00:01:21 ===========read db of one=======
[xorm] [info]  2019/11/12 00:01:21.884865 [SQL] SELECT `id`, `name`, `age` FROM `user` WHERE (id = ?) LIMIT 1 []interface {}{1} - took: 2.563032ms
2019/11/12 00:01:21 read one db,get id = 1 of userInfo:  true <nil>
2019/11/12 00:01:21 =========slave two db====
[xorm] [info]  2019/11/12 00:01:21.894919 [SQL] SELECT `id`, `name`, `age` FROM `user` WHERE (id = ?) LIMIT 1 []interface {}{2} - took: 9.691442ms
2019/11/12 00:01:21 read two db user:  true <nil>
2019/11/12 00:01:21 =======engine select=========
[xorm] [info]  2019/11/12 00:01:21.895281 [SQL] SELECT `id`, `name`, `age` FROM `user` WHERE (id = ?) LIMIT 1 []interface {}{3}
2019/11/12 00:01:21 true <nil>
2019/11/12 00:01:21 &{3 xiaoxiao 12}
[xorm] [info]  2019/11/12 00:01:21.896021 [SQL] INSERT INTO `user` (`name`,`age`) VALUES (?, ?) []interface {}{"xiaoxiao", 12}
2019/11/12 00:01:21 affected num:  1
2019/11/12 00:01:21 insert id:  104
2019/11/12 00:01:21 err:  <nil>
2019/11/12 00:01:21 get on slave to query
[xorm] [info]  2019/11/12 00:01:21.953948 [SQL] SELECT `id`, `name`, `age` FROM `user` WHERE (id = ?) LIMIT 1 []interface {}{4} - took: 917.993µs
2019/11/12 00:01:21 true <nil>
--- PASS: TestGxorm (0.12s)
PASS
ok      github.com/daheige/thinkgo/gxorm        0.128s
*/
