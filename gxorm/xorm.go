package gxorm

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
	Loc          string //时区字符串 Local,PRC
	SqlCmd       bool   //sql语句是否输出到终端,true输出到终端
	UsePool      bool   //当前db实例是否采用db连接池,默认不采用，如采用请求配置该参数
	ShowExecTime bool   //是否打印sql执行时间
}

//每个数据库连接pool就是一个db引擎
var engineMap = map[string]*xorm.Engine{}

func (conf *DbConf) InitEngine() (*xorm.Engine, error) {
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
		return nil, err
	}

	if conf.SqlCmd {
		db.ShowSQL(true) //控制台打印出sql
	}

	if conf.ShowExecTime {
		db.ShowExecTime(true)
	}

	//设置连接池
	if conf.UsePool {
		db.SetMaxIdleConns(conf.MaxIdleConns) //设置连接池的空闲数大小
		db.SetMaxOpenConns(conf.MaxOpenConns) //设置最大打开连接数
	}

	return db, nil
}

// SetEngineName 给当前数据库指定engineName
// 一般用在多个db 数据库连接引擎的时候，可以给当前的db engine设置一个name
// 这样业务上游层，就可以通过 GetEngine(name)获得当前db engine
func (conf *DbConf) SetEngineName(name string) error {
	if name == "" {
		return errors.New("current engineGroup name is empty!")
	}

	//初始化db 句柄
	db, err := conf.InitEngine()
	if err != nil {
		return errors.New("current " + name + " db engine init error: " + err.Error())
	}

	engineMap[name] = db
	return nil
}

// ShortConnect 短连接设置，一般用于短连接服务的数据库句柄
func (conf *DbConf) ShortConnect() (*xorm.Engine, error) {
	conf.UsePool = false
	return conf.InitEngine()
}

//从db pool获取一个数据库连接句柄
//根据数据库连接句柄name获取指定的连接句柄
func GetEngine(name string) (*xorm.Engine, error) {
	if _, ok := engineMap[name]; ok {
		return engineMap[name], nil
	}

	return nil, errors.New("get db obj failed!")
}

// CloseAllDb 由于xorm db.Close()是关闭当前连接，一般建议如下函数放在main/init关闭连接就可以
func CloseAllDb() {
	for name, db := range engineMap {
		if err := db.Close(); err != nil {
			log.Println("close db error: ", err.Error())
			continue
		}

		delete(engineMap, name) //销毁连接句柄标识
	}
}

// CloseDbByName 关闭指定name的db engine
func CloseDbByName(name string) error {
	if _, ok := engineMap[name]; ok {
		if err := engineMap[name].Close(); err != nil {
			log.Println("close db error: ", err.Error())
			return err
		}

		delete(engineMap, name)
	}

	return errors.New("current db engine not exist")
}

//======================读写分离设置==================
// NewEngineGroup 设置读写分离db engine
// slaveEngine可以多个
// 返回读写分离的db engine
func NewEngineGroup(masterEngine *xorm.Engine, slave1Engine ...*xorm.Engine) (*xorm.EngineGroup, error) {
	engineGroup, err := xorm.NewEngineGroup(masterEngine, slave1Engine)
	if err != nil {
		return nil, err
	}

	return engineGroup, nil
}
