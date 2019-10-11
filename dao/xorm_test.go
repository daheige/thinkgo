package dao

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

func TestDao(t *testing.T) {
	var e *xorm.Engine
	log.Println(e == nil)

	dbconf := &DbConf{
		Ip:           "127.0.0.1",
		Port:         3306,
		User:         "root",
		Password:     "root",
		Database:     "test",
		MaxIdleConns: 10,
		MaxOpenConns: 100,
		ParseTime:    true,
		SqlCmd:       true,
		ShowExecTime: false,
	}

	dbconf.SetEngine() //设置数据库连接对象，并非真正连接，只有在用的时候才会建立连接
	dbconf.SetEngineName("default")

	db, err := dbconf.Db()
	if db == nil || err != nil {
		log.Println("db error")
		return
	}

	defer db.Close()

	user := &myUser{}
	log.Println("====user1====")
	has, err := db.Where("id = ?", 1).Get(user)
	log.Println(has, err)

	log.Println("user: ", user)
	log.Println(user.Id, user.Name)

	//测试读写分离
	readConf := &DbConf{
		Ip:           "127.0.0.1",
		Port:         3306,
		User:         "test1",
		Password:     "1234",
		Database:     "test",
		MaxIdleConns: 10,
		MaxOpenConns: 100,
		ParseTime:    true,
		SqlCmd:       true,
		ShowExecTime: true,
	}

	readConf.SetEngine()
	// readConf.SetEngineName("readEngine") //为每个db设置一个engine name

	readEngine, err := readConf.Db()
	if err != nil {
		log.Println("set read db engine error: ", err.Error())
		return
	}

	defer readEngine.Close()

	userInfo := &myUser{}
	has, err = readEngine.Where("id = ?", 1).Get(userInfo)
	log.Println("read one db,get id = 1 of userInfo: ", has, err)

	//设置第二个读的实例
	readConf2 := &DbConf{
		Ip:           "127.0.0.1",
		Port:         3306,
		User:         "test2",
		Password:     "1234",
		Database:     "test",
		MaxIdleConns: 10,
		MaxOpenConns: 100,
		ParseTime:    true,
		SqlCmd:       true,
		ShowExecTime: true,
	}

	readConf2.SetEngine()
	// readConf2.SetEngineName("readEngine2") //为每个db设置一个engine name
	readEngine2, err := readConf2.Db()
	if err != nil {
		log.Println("set read db engine error: ", err.Error())
		return
	}

	defer readConf2.Close()

	userInfo2 := &myUser{}
	has, err = readEngine2.Where("id = ?", 2).Get(userInfo2)
	log.Println("read two db user: ", has, err)

	//设置读写分离的引擎句柄
	// engineGroup, err := SetEngineGroup(db, readEngine)
	// engineGroup, err := SetEngineGroup(db, readEngine2)
	engineGroup, err := SetEngineGroup(db, readEngine, readEngine2)
	if err != nil {
		log.Println("set db engineGroup error: ", err.Error())
		return
	}

	//defer engineGroup.Close() //关闭读写分离的连接对象

	//设置读写分离连接对象，并非真正建立连接
	SetEngineGroupName("dbGroup", engineGroup) //为每个db设置一个engine name

	log.Println("=======engine select=========")
	user2 := &myUser{}
	has, err = engineGroup.Where("id = ?", 3).Get(user2)
	log.Println(has, err)
	log.Println(user2)

	//通过辅助函数获取读写分离连接对象
	db2, err := GetEngineGroup("dbGroup")
	if err != nil {
		log.Println(err)
		return
	}

	user3 := &myUser{}
	has, err = db2.Where("id = ?", 4).Get(user3)
	log.Println("====user3====")
	log.Println(has, err)
	log.Println(user3)

	//采用读写分离实现数据插入
	user4 := &myUser{
		Name: "xiaoxiao",
		Age:  12,
	}

	affectedNum, err := db2.InsertOne(user4) //插入单条数据，多条数据请用Insert(user3,user4,user5)
	log.Println("affected num: ", affectedNum)
	log.Println("insert id: ", user4.Id)
	log.Println("err: ", err)

	log.Println("get on slave to query")
	user5 := &myUser{}
	log.Println(engineGroup.Slave().Where("id = ?", 4).Get(user5))
}

/**
$ go test -v
=== RUN   TestDao
2019/10/11 22:42:23 true
2019/10/11 22:42:23 ====user1====
[xorm] [info]  2019/10/11 22:42:23.496756 [SQL] SELECT `id`, `name`, `age` FROM `user` WHERE (id = ?) LIMIT 1 []interface {}{1}
2019/10/11 22:42:23 true <nil>
2019/10/11 22:42:23 user:  &{1 xiaoxiao 12}
2019/10/11 22:42:23 1 xiaoxiao
[xorm] [info]  2019/10/11 22:42:23.499520 [SQL] SELECT `id`, `name`, `age` FROM `user` WHERE (id = ?) LIMIT 1 []interface {}{1} - took: 983.163µs
2019/10/11 22:42:23 read one db,get id = 1 of userInfo:  true <nil>
[xorm] [info]  2019/10/11 22:42:23.501532 [SQL] SELECT `id`, `name`, `age` FROM `user` WHERE (id = ?) LIMIT 1 []interface {}{2} - took: 1.247577ms
2019/10/11 22:42:23 read two db user:  true <nil>
2019/10/11 22:42:23 =======engine select=========
[xorm] [info]  2019/10/11 22:42:23.501695 [SQL] SELECT `id`, `name`, `age` FROM `user` WHERE (id = ?) LIMIT 1 []interface {}{3}
2019/10/11 22:42:23 true <nil>
2019/10/11 22:42:23 &{3 xiaoxiao 12}
[xorm] [info]  2019/10/11 22:42:23.502124 [SQL] SELECT `id`, `name`, `age` FROM `user` WHERE (id = ?) LIMIT 1 []interface {}{4}
2019/10/11 22:42:23 ====user3====
2019/10/11 22:42:23 true <nil>
2019/10/11 22:42:23 &{4 xiaoxiao 12}
[xorm] [info]  2019/10/11 22:42:23.502672 [SQL] INSERT INTO `user` (`name`, `age`) VALUES (?, ?) []interface {}{"xiaoxiao", 12}
2019/10/11 22:42:23 affected num:  1
2019/10/11 22:42:23 insert id:  96
2019/10/11 22:42:23 err:  <nil>
2019/10/11 22:42:23 get on slave to query
[xorm] [info]  2019/10/11 22:42:23.543954 [SQL] SELECT `id`, `name`, `age` FROM `user` WHERE (id = ?) LIMIT 1 []interface {}{4} - took: 805.287µs
2019/10/11 22:42:23 true <nil>
--- PASS: TestDao (0.05s)
PASS
ok  	github.com/daheige/thinkgo/dao	0.052s
*/
