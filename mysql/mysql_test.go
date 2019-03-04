package mysql

import (
	"log"
	"sync"
	"testing"
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

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			if err != nil {
				log.Println("get db error: ", err.Error())
			}

			user := &myUser{}
			db.Where("name = ?", "hello").First(user)
			log.Println(user)
		}()
	}
}

/** go test -v
 * 2018/10/27 12:31:48 test success
--- PASS: TestGorm (1.38s)
PASS
ok      github.com/daheige/thinkgo/gorm/mysql  1.392s

*/
