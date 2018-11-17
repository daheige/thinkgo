/**
* gorm mysql封装，支持多个数据库实例化为连接池对象
* 结合了xorm思想，将每个数据库对象作为一个数据库引擎句柄
* xorm设计思想：在xorm里面，可以同时存在多个Orm引擎
* 一个Orm引擎称为Engine，一个Engine一般只对应一个数据库
* 因此,可以将gorm的每个数据库连接句柄，可以作为一个Engine来进行处理
 */
package mysql

import (
	"errors"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

//每个数据库连接pool就是一个db引擎
var engineMap = map[string]*gorm.DB{}

//mysql连接信息
//parseTime=true changes the output type of DATE and DATETIME
//values to time.Time instead of []byte / string
//The date or datetime like 0000-00-00 00:00:00 is converted
//into zero value of time.Time.
type DbConf struct {
	Ip           string
	Port         int
	User         string
	Password     string
	Database     string
	Charset      string //字符集 utf8mb4 支持表情符号
	Collation    string //整理字符集 utf8mb4_unicode_ci
	MaxIdleConns int    //空闲pool个数
	MaxOpenConns int    //最大open connection个数
	ParseTime    bool
	Loc          string   //时区字符串 Local,PRC
	engineName   string   //当前数据库连接句柄标识
	dbObj        *gorm.DB //当前数据库连接句柄
}

//给当前数据库指定engineName
func (conf *DbConf) SetEngineName(name string) {
	conf.engineName = name
}

//给当前数据库指定一个engine name,并建立连接池
func (conf *DbConf) SetDbPool() error {
	if conf.engineName == "" {
		panic("name must be not null")
	}

	err := conf.initDb()
	if err != nil {
		return errors.New("set dbEngine failed")
	}

	engineMap[conf.engineName] = conf.dbObj
	return nil
}

//关闭当前数据库连接
// 一般建议，将当前db engine close函数放在main/init关闭连接就可以
func (conf *DbConf) Close() error {
	if db, ok := engineMap[conf.engineName]; ok {
		if err := db.Close(); err != nil {
			log.Println("close db error: ", err.Error())
			return err
		}

		return nil
	}

	return errors.New("current db engine not exist")
}

func (conf *DbConf) Db() *gorm.DB {
	return conf.dbObj
}

//建立db连接对象，并非真正连接db,只有在用的时候才会建立db连接
func (conf *DbConf) initDb() error {
	if conf.Ip == "" {
		conf.Ip = "127.0.0.1"
	}

	if conf.Port == 0 {
		conf.Port = 3306
	}

	if conf.Charset == "" {
		conf.Charset = "utf8mb4"
	}

	if conf.Collation == "" {
		conf.Collation = "utf8mb4_unicode_ci"
	}

	if conf.Loc == "" {
		conf.Loc = "Local"
	}

	//连接gorm.DB实例对象，并非立即连接db,用的时候才会真正的建立连接
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&collation=%s&parseTime=%v&loc=%s",
		conf.User, conf.Password, conf.Ip, conf.Port, conf.Database,
		conf.Charset, conf.Collation, conf.ParseTime, conf.Loc))
	if err != nil {
		return errors.New("set gorm.DB failed")
	}

	//设置连接池
	db.DB().SetMaxIdleConns(conf.MaxIdleConns)
	db.DB().SetMaxOpenConns(conf.MaxOpenConns)

	conf.dbObj = db
	return nil
}

// ========================辅助函数===============
//从db pool获取一个数据库连接句柄
//根据数据库连接句柄name获取指定的连接句柄
func GetDbObj(name string) (*gorm.DB, error) {
	if _, ok := engineMap[name]; ok {
		return engineMap[name], nil
	}

	return nil, errors.New("get db obj failed")
}

//由于gorm db.Close()是关闭当前连接，一般建议如下函数放在main/init关闭连接就可以
func CloseAllDb() {
	for name, db := range engineMap {
		if err := db.Close(); err != nil {
			log.Println("close db error: ", err.Error())
			continue
		}

		delete(engineMap, name) //销毁连接句柄标识
	}
}

//关闭指定name的db engine
func CloseDbByName(name string) error {
	if _, ok := engineMap[name]; ok {
		if err := engineMap[name].Close(); err != nil {
			log.Println("close db error: ", err.Error())
			return err
		}

		delete(engineMap, name) //销毁连接句柄标识
	}

	return errors.New("current dbObj not exist")
}