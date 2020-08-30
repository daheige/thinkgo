/** Package mysql of gorm library.
* gorm mysql封装，支持多个数据库实例化为连接池对象
* 结合了xorm思想，将每个数据库对象作为一个数据库引擎句柄
* xorm设计思想：在xorm里面，可以同时存在多个Orm引擎
* 一个Orm引擎称为Engine，一个Engine一般只对应一个数据库
* 因此,可以将gorm的每个数据库连接句柄，可以作为一个Engine来进行处理
* 容易踩坑的地方：
*	对于golang的官方sql引擎，sql.open并非立即连接db,用的时候才会真正的建立连接
*	但是gorm.Open在设置完db对象后，还发送了一个Ping操作，判断连接是否连接上去
*	对于短连接的话，建议用完就调用db.Close()方法释放db连接资源
*	对于长连接服务，一般建议在main/init中关闭连接就可以
*	具体可以看gorm/main.go源码85行
* 对于gorm实现读写分离:
*	可以实例化master,slaves实例，对于curd用不同的句柄就可以
 */
package mysql

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// engineMap 每个数据库连接pool就是一个db引擎
var engineMap = map[string]*gorm.DB{}

// DbConf mysql连接信息
// parseTime=true changes the output type of DATE and DATETIME
// values to time.Time instead of []byte / string
// The date or datetime like 0000-00-00 00:00:00 is converted
// into zero value of time.Time.
type DbConf struct {
	Ip        string
	Port      int // 默认3306
	User      string
	Password  string
	Database  string
	Charset   string // 字符集 utf8mb4 支持表情符号
	Collation string // 整理字符集 utf8mb4_unicode_ci

	UsePool      bool // 当前db实例是否采用db连接池,默认不采用，如采用请求配置该参数
	MaxIdleConns int  // 空闲pool个数
	MaxOpenConns int  // 最大open connection个数

	// sets the maximum amount of time a connection may be reused.
	// 设置连接可以重用的最大时间
	// 给db设置一个超时时间，时间小于数据库的超时时间
	MaxLifetime int64 // 数据库超时时间，单位s

	// 连接超时/读取超时/写入超时设置
	Timeout      time.Duration // Dial timeout
	ReadTimeout  time.Duration // I/O read timeout
	WriteTimeout time.Duration // I/O write timeout

	ParseTime  bool     // 格式化时间类型
	Loc        string   // 时区字符串 Local,PRC
	engineName string   // 当前数据库连接句柄标识
	dbObj      *gorm.DB // 当前数据库连接句柄

	ShowSql bool           // sql语句是否输出
	Logger  gorm.LogWriter // sql输出logger句柄接口
}

// SetEngineName 给当前数据库指定engineName
func (conf *DbConf) SetEngineName(name string) error {
	if name == "" {
		return errors.New("current engine name is empty!")
	}

	if conf.dbObj == nil {
		return errors.New("current " + name + " db engine not be initDb")
	}

	conf.engineName = name
	engineMap[conf.engineName] = conf.dbObj

	return nil
}

// SetDbObj 创建当前数据库db对象，并非连接，在使用的时候才会真正建立db连接
// 为兼容之前的版本，这里新增SetDb创建db对象
func (conf *DbConf) SetDbObj() error {
	err := conf.initDb()
	if err != nil {
		log.Println("set db engine error: ", err)
		return err
	}

	return nil
}

// SetDbPool 设置db pool连接池
func (conf *DbConf) SetDbPool() error {
	conf.UsePool = true
	return conf.SetDbObj()
}

// ShortConnect 建立短连接，用完需要调用Close()进行关闭连接，释放资源，否则就会出现too many connection
func (conf *DbConf) ShortConnect() error {
	conf.UsePool = false
	err := conf.initDb()
	if err != nil {
		log.Println("set db engine error: ", err)
		return err
	}

	return nil
}

// Close 关闭当前数据库连接
// 一般建议，将当前db engine close函数放在main/init关闭连接就可以
func (conf *DbConf) Close() error {
	if conf.dbObj == nil {
		return nil
	}

	if err := conf.Db().Close(); err != nil {
		log.Println("close db error: ", err.Error())
		return err
	}

	// 把连接句柄对象从map中删除
	delete(engineMap, conf.engineName)

	return nil
}

// Db 返回当前db对象
func (conf *DbConf) Db() *gorm.DB {
	return conf.dbObj
}

// initDb 建立db连接句柄
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

	if conf.Timeout == 0 {
		conf.Timeout = 10 * time.Second
	}

	if conf.WriteTimeout == 0 {
		conf.WriteTimeout = 5 * time.Second
	}

	if conf.ReadTimeout == 0 {
		conf.ReadTimeout = 5 * time.Second
	}

	// mysql connection time loc.
	loc, err := time.LoadLocation(conf.Loc)
	if err != nil {
		return err
	}

	// mysql config
	mysqlConf := mysql.Config{
		User:   conf.User,
		Passwd: conf.Password,
		Net:    "tcp",
		Addr:   fmt.Sprintf("%s:%d", conf.Ip, conf.Port),
		DBName: conf.Database,
		// Connection parameters
		Params: map[string]string{
			"charset": conf.Charset,
		},
		Collation:            conf.Collation,
		Loc:                  loc,               // Location for time.Time values
		Timeout:              conf.Timeout,      // Dial timeout
		ReadTimeout:          conf.ReadTimeout,  // I/O read timeout
		WriteTimeout:         conf.WriteTimeout, // I/O write timeout
		AllowNativePasswords: true,              // Allows the native password authentication method
		ParseTime:            conf.ParseTime,    // Parse time values to time.Time
	}

	// 对于golang的官方sql引擎，sql.open并非立即连接db,用的时候才会真正的建立连接
	// 但是gorm.Open在设置完db对象后，还发送了一个Ping操作，判断连接是否连接上去
	// 具体可以看gorm/main.go源码Open方法
	db, err := gorm.Open("mysql", mysqlConf.FormatDSN())
	if err != nil { // 数据库连接错误
		log.Println("open mysql connection error: ", err)

		return err
	}

	if conf.ShowSql {
		db.LogMode(true)
		if conf.Logger == nil {
			conf.Logger = log.New(os.Stdout, "\r\n", 0)
		}

		db.SetLogger(gorm.Logger{
			LogWriter: conf.Logger,
		})
	}

	// 设置连接池
	if conf.UsePool {
		db.DB().SetMaxIdleConns(conf.MaxIdleConns)
		db.DB().SetMaxOpenConns(conf.MaxOpenConns)
	}

	// 设置连接可以重用的最大时间
	// 给db设置一个超时时间，时间小于数据库的超时时间
	if conf.MaxLifetime > 0 {
		db.DB().SetConnMaxLifetime(time.Duration(conf.MaxLifetime) * time.Second)
	}

	conf.dbObj = db

	return nil
}

// ========================辅助函数===============
// GetDbObj 从db pool获取一个数据库连接句柄
// 根据数据库连接句柄name获取指定的连接句柄
func GetDbObj(name string) (*gorm.DB, error) {
	if _, ok := engineMap[name]; ok {
		return engineMap[name], nil
	}

	return nil, errors.New("get db obj failed")
}

// CloseAllDb 由于gorm db.Close()是关闭当前连接，一般建议如下函数放在main/init关闭连接就可以
func CloseAllDb() {
	for name, db := range engineMap {
		if err := db.Close(); err != nil {
			log.Println("close db error: ", err.Error())
			continue
		}

		delete(engineMap, name) // 销毁连接句柄标识
	}
}

// CloseDbByName 关闭指定name的db engine
func CloseDbByName(name string) error {
	if _, ok := engineMap[name]; ok {
		if err := engineMap[name].Close(); err != nil {
			log.Println("close db error: ", err.Error())
			return err
		}

		delete(engineMap, name) // 销毁连接句柄标识
	}

	return errors.New("current dbObj not exist")
}
