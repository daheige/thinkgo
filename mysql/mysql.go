/**
* gorm mysql封装，支持多个数据库实例化为连接池对象
* 结合了xorm思想，将每个数据库对象作为一个数据库引擎句柄
* xorm设计思想：在xorm里面，可以同时存在多个Orm引擎
* 一个Orm引擎称为Engine，一个Engine一般只对应一个数据库
* 因此,可以将gorm的每个数据库连接句柄，可以作为一个Engine来进行处理
* 容易踩坑的地方：
	对于golang的官方sql引擎，sql.open并非立即连接db,用的时候才会真正的建立连接
	但是gorm.Open在设置完db对象后，还发送了一个Ping操作，判断连接是否连接上去
	对于短连接的话，建议用完就调用db.Close()方法释放db连接资源
	对于长连接服务，一般建议在main/init中关闭连接就可以
	具体可以看gorm/main.go源码85行
*/
package mysql

import (
	"errors"
	"fmt"
	"log"
	"os"

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
	SqlCmd       bool     //sql语句是否输出到终端,true输出到终端
	UsePool      bool     //当前db实例是否采用db连接池,默认不采用，如采用请求配置该参数

	//默认采用gorm log输出sql语句到终端
	//当 defaultLogType=true,就采用自定义的stdout输出sql日志
	defaultLogType bool
}

//给当前数据库指定engineName
func (conf *DbConf) SetEngineName(name string) {
	if name == "" {
		panic(fmt.Sprintf("current %s engine name is empty!", name))
	}

	if conf.dbObj == nil {
		panic(fmt.Sprintf("current %s db engine not initDb", name))
	}

	conf.engineName = name
	engineMap[conf.engineName] = conf.dbObj
}

//创建当前数据库db对象，并非连接，在使用的时候才会真正建立db连接
//为兼容之前的版本，这里新增SetDb创建db对象
func (conf *DbConf) SetDbObj() error {
	err := conf.initDb()
	if err != nil {
		return errors.New("set dbEngine failed")
	}

	return nil
}

//设置日志自定义os.Stdout输出格式
func (conf *DbConf) setLogType(flag bool) {
	conf.defaultLogType = flag
}

//设置db pool连接池
func (conf *DbConf) SetDbPool() error {
	conf.UsePool = true
	return conf.SetDbObj()
}

//建立短连接，用完需要调用Close()进行关闭连接，释放资源，否则就会出现too many connection
func (conf *DbConf) ShortConnect() error {
	conf.UsePool = false
	err := conf.initDb()
	if err != nil {
		return errors.New("set dbEngine failed")
	}

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

//返回当前db对象
func (conf *DbConf) Db() *gorm.DB {
	return conf.dbObj
}

//建立db连接句柄
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

	//对于golang的官方sql引擎，sql.open并非立即连接db,用的时候才会真正的建立连接
	//但是gorm.Open在设置完db对象后，还发送了一个Ping操作，判断连接是否连接上去
	//具体可以看gorm/main.go源码85行
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&collation=%s&parseTime=%v&loc=%s",
		conf.User, conf.Password, conf.Ip, conf.Port, conf.Database,
		conf.Charset, conf.Collation, conf.ParseTime, conf.Loc))
	if err != nil {
		return errors.New("set gorm.DB failed")
	}

	//将sql打印到终端
	if conf.SqlCmd {
		db.LogMode(true)
		if conf.defaultLogType { //采用os.Stdout输出日志格式
			db.SetLogger(log.New(os.Stdout, "\n", log.LstdFlags))
		}
	}

	//设置连接池
	if conf.UsePool {
		db.DB().SetMaxIdleConns(conf.MaxIdleConns)
		db.DB().SetMaxOpenConns(conf.MaxOpenConns)
	}

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
