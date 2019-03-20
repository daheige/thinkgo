package mysql

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"testing"

	"github.com/jinzhu/gorm"
)

/**
* sql:
* CREATE DATABASE IF NOT EXISTS test default charset utf8mb4;
* create table user (id int primary key auto_increment,name varchar(200)) engine=innodb;
* 模拟数据插入
* mysql> insert into user (name) values("xiaoming");
   Query OK, 1 row affected (0.11 sec)

   mysql> insert into user (name) values("hello");
   Query OK, 1 row affected (0.04 sec)
*/
type myUser struct {
	ID   uint   `gorm:"primary_key"`
	Name string `gorm:"type:varchar(200)"`
}

func (myUser) TableName() string {
	return "user"
}

func TestGorm(t *testing.T) {
	dbconf := &DbConf{
		Ip:           "127.0.0.1",
		Port:         3306,
		User:         "root",
		Password:     "1234",
		Database:     "test",
		MaxIdleConns: 10,
		MaxOpenConns: 100,
		ParseTime:    true,
		SqlCmd:       true,
	}

	dbconf.setLogType(true)

	//设置db engine name
	dbconf.SetEngineName("default")
	dbconf.SetDbPool()   //建立db连接池
	defer dbconf.Close() //关闭当前连接
	// defer CloseAllDb() //关闭所有的连接句柄

	db, err := GetDbObj("default")
	if err != nil {
		t.Log("get db error: ", err.Error())
	}

	user := &myUser{}
	db.Where("name = ?", "hello").First(user)
	log.Println(user)

	var wg sync.WaitGroup
	testFind(&wg)

	wg.Wait()
	log.Println("test success")
}

func testFind(wg *sync.WaitGroup) {
	db, err := GetDbObj("default")
	defer db.Close() //当我们在这里进行关了db close相当于断开连接

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// db, err := GetDbObj("default")

			if err != nil {
				log.Println("get db error: ", err.Error())
				return
			}

			user := &myUser{}
			db.Where("name = ?", "hello").First(user)
			log.Println(user)
		}()
	}
}

func TestShortConnect(t *testing.T) {
	getDb := func() (*gorm.DB, error) {
		conf := &DbConf{
			Ip:           "127.0.0.1",
			Port:         3306,
			User:         "root",
			Password:     "1234",
			Database:     "test",
			MaxIdleConns: 10,
			MaxOpenConns: 100,
			ParseTime:    true,
			SqlCmd:       true,
			Charset:      "utf8mb4",
			Collation:    "utf8mb4_unicode_ci",
			Loc:          "Local",
		}

		//连接gorm.DB实例对象，并非立即连接db,用的时候才会真正的建立连接
		db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&collation=%s&parseTime=%v&loc=%s",
			conf.User, conf.Password, conf.Ip, conf.Port, conf.Database,
			conf.Charset, conf.Collation, conf.ParseTime, conf.Loc))
		if err != nil {
			return nil, errors.New("set gorm.DB failed")
		}

		return db, nil
	}

	//这里我设置了db max_connections最大连接为1000
	var wg sync.WaitGroup
	var maxConnections = 1024
	// var maxConnections = 2000
	for i := 0; i < maxConnections; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			db, err := getDb()
			defer db.Close()

			if err != nil {
				log.Println("get db error: ", err.Error())
				return
			}

			user := &myUser{}
			db.Where("name = ?", "hello").First(user)
			log.Println(user)
		}()
	}

	wg.Wait()
	log.Println("test success")
}

/** go test -v -test.run TestGorm
采用长连接测试
--- PASS: TestGorm (1.35s)
ok  	github.com/daheige/thinkgo/mysql	1.365s
采用短连接方式测试
	go test -v -test.run TestShortConnect
--- PASS: TestShortConnect (1.22s)
PASS
ok  	github.com/daheige/thinkgo/mysql	1.231s

当我们把maxConnections 调到2000后
$ go test -v -test.run TestShortConnect
=== RUN   TestShortConnect
2019/03/20 15:15:06 get db error:  set gorm.DB failed
panic: runtime error: invalid memory address or nil pointer dereference
[signal SIGSEGV: segmentation violation code=0x1 addr=0x88 pc=0x6be466]

goroutine 1401 [running]:
就会出现panic
*/
