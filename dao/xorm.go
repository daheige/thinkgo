package dao

import (
	"errors"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/go-xorm/xorm"
)

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
	Loc          string       //时区字符串 Local,PRC
	engineName   string   //当前数据库连接句柄标识
	dbObj        *xorm.Engine //当前数据库连接句柄
	SqlCmd       bool         //sql语句是否输出到终端,true输出到终端
	UsePool      bool         //当前db实例是否采用db连接池,默认不采用，如采用请求配置该参数
}

//每个数据库连接pool就是一个db引擎
var engineMap = map[string]*xorm.Engine{}

//读写分离的引擎组
var engineGroupMap = map[string]*xorm.EngineGroup{}

//返回当前db实例对象，并非立即连接db,用的时候才会真正的建立连接
func (conf *DbConf) Db() (*xorm.Engine, error) {
	if conf.dbObj == nil {
		return nil, errors.New("current db engine not initDb")
	}

	return conf.dbObj, nil
}

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

	//连接实例对象，并非立即连接db,用的时候才会真正的建立连接
	db, err := xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&collation=%s&parseTime=%v&loc=%s",
		conf.User, conf.Password, conf.Ip, conf.Port, conf.Database,
		conf.Charset, conf.Collation, conf.ParseTime, conf.Loc))
	if err != nil {
		return err
	}

	if conf.SqlCmd {
		db.ShowSQL(true) //控制台打印出sql
	}

	//设置连接池
	if conf.UsePool {
		db.SetMaxIdleConns(conf.MaxIdleConns) //设置连接池的空闲数大小
		db.SetMaxOpenConns(conf.MaxOpenConns) //设置最大打开连接数
	}

	//当前数据库连接对象
	conf.dbObj = db
	return nil
}

//创建当前数据库db对象，并非连接，在使用的时候才会真正建立db连接
func (conf *DbConf) SetEngine() error {
	err := conf.initDb()
	if err != nil {
		return errors.New("set db engine failed")
	}

	return nil
}

//给当前数据库指定engineName
func (conf *DbConf) SetEngineName(name string) {
	if name == "" {
		panic(fmt.Sprintf("current %s engineGroup name is empty!", name))
	}

	if conf.dbObj == nil {
		panic(fmt.Sprintf("current %s db engine not initDb", name))
	}

	conf.engineName = name
	engineMap[name] = conf.dbObj
}

//短连接设置，一般用于短连接服务的数据库句柄
func (conf *DbConf) ShortConnect() error {
	conf.UsePool = false
	err := conf.initDb()
	if err != nil {
		return errors.New("set dbEngine failed")
	}

	return nil
}

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


//======================读写分离设置==================
//读写分离，masterEngine,slaveEngine可以多个
func SetEngineGroup(masterEngine *xorm.Engine, slave1Engine ...*xorm.Engine) (*xorm.EngineGroup, error) {
	engineGroup, err := xorm.NewEngineGroup(masterEngine, slave1Engine)
	if err != nil {
		return nil, err
	}

	return engineGroup, nil
}

//给读写分离的dbGroup设置name，一般用于业务上游层调度
func SetEngineGroupName(name string, engineGroup *xorm.EngineGroup) {
	if name == "" {
		panic(fmt.Sprintf("current %s engineGroup name is empty", name))
	}

	if engineGroup == nil {
		panic(fmt.Sprintf("current %s engineGroup is nil", name))
	}

	//设置读写分离的句柄名称
	engineGroupMap[name] = engineGroup
}

//通过name获取读写分离的句柄对象,并非真正建立连接，只有在使用的时候才会建立连接
func GetEngineGroup(name string) (*xorm.EngineGroup, error) {
	if _, ok := engineGroupMap[name]; ok {
		return engineGroupMap[name], nil
	}

	return nil, errors.New(fmt.Sprintf("current %s engineGroup not exist!", name))
}

//=============================读写分离，关闭连接对象，辅助函数====================
//关闭所有的读写分离连接句柄对象
func CloseAllEngineGroup() {
	for name, db := range engineGroupMap {
		if err := db.Close(); err != nil {
			log.Println("close db error: ", err.Error())
			continue
		}

		delete(engineGroupMap, name) //销毁连接句柄标识
	}
}

//根据读写分离的句柄名称关闭连接对象
func CloseEngineGroupByName(name string) error {
	if _, ok := engineGroupMap[name]; ok {
		if err := engineGroupMap[name].Close(); err != nil {
			log.Println("close db error: ", err.Error())
			return err
		}

		delete(engineGroupMap, name) //销毁连接句柄标识
	}

	return errors.New("current enginGroup not exist")
}

//===================对于非读写分离模式下，单个数据库引擎，辅助函数===========

//从db pool获取一个数据库连接句柄
//根据数据库连接句柄name获取指定的连接句柄
func GetEngine(name string) (*xorm.Engine, error) {
	if _, ok := engineMap[name]; ok {
		return engineMap[name], nil
	}

	return nil, errors.New("get db obj failed!")
}

//由于xorm db.Close()是关闭当前连接，一般建议如下函数放在main/init关闭连接就可以
func CloseAllEngine() {
	for name, db := range engineMap {
		if err := db.Close(); err != nil {
			log.Println("close db error: ", err.Error())
			continue
		}

		delete(engineMap, name) //销毁连接句柄标识
	}
}

//关闭指定name的db engine
func CloseEngineByName(name string) error {
	if _, ok := engineMap[name]; ok {
		if err := engineMap[name].Close(); err != nil {
			log.Println("close db error: ", err.Error())
			return err
		}

		delete(engineMap, name) //销毁连接句柄标识
	}

	return errors.New("current dbObj not exist")
}
