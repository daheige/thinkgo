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
* 由于gorm自己对mysql做了一次包裹，所以重命名处理
* gMysql "gorm.io/driver/mysql"
* gorm v2版本仓库地址：https://github.com/go-gorm/gorm
 */
package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
	gMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	hasInit    bool     // 是否调用了 InitInstance()进行初始化db

	ShowSql bool // sql语句是否输出

	// sql输出logger句柄接口
	// logger.Writer 接口需要实现Printf(string, ...interface{}) 方法
	// 具体可以看gorm v2 logger包源码
	// https://github.com/go-gorm/gorm
	Logger logger.Writer

	// gorm v2版本新增参数
	gMysqlConfig gMysql.Config // gorm v2新增参数gMysql.Config
	gormConfig   gorm.Config   // gorm v2新增参数gorm.Config
	LoggerConfig logger.Config // gorm v2新增参数logger.Config
}

// SetEngineName 给当前数据库指定engineName
func (conf *DbConf) SetEngineName(name string) error {
	if name == "" {
		return errors.New("current engine name is empty")
	}

	if !conf.hasInit {
		return errors.New("current " + name + " db engine must be InitInstance")
	}

	conf.engineName = name
	engineMap[conf.engineName] = conf.dbObj

	return nil
}

// SetDbPool 设置db pool连接池
func (conf *DbConf) SetDbPool() error {
	conf.UsePool = true
	return conf.InitInstance()
}

// ShortConnect 建立短连接，用完需要调用Close()进行关闭连接，释放资源
// 否则就会出现too many connection
func (conf *DbConf) ShortConnect() error {
	conf.UsePool = false
	return conf.InitInstance()
}

// Close 关闭当前数据库连接
// 一般建议，将当前db engine close函数放在main/init关闭连接就可以
func (conf *DbConf) Close() error {
	if conf.dbObj == nil {
		return nil
	}

	db, err := conf.SqlDB()
	if err != nil {
		log.Println("get db instance error: ", err.Error())
		return err
	}

	err = db.Close()
	if err != nil {
		log.Println("close db instance error: ", err.Error())
		return err
	}

	if conf.engineName != "" {
		// 把连接句柄对象从map中删除
		delete(engineMap, conf.engineName)
	}

	return nil
}

// Db 返回当前db对象
func (conf *DbConf) Db() *gorm.DB {
	return conf.dbObj
}

// SqlDB 返回sql DB
func (conf *DbConf) SqlDB() (*sql.DB, error) {
	return conf.dbObj.DB()
}

// DSN 设置mysql dsn
// mysql charset查看
// mysql> show character set where charset="utf8mb4";
// +---------+---------------+--------------------+--------+
// | Charset | Description   | Default collation  | Maxlen |
// +---------+---------------+--------------------+--------+
// | utf8mb4 | UTF-8 Unicode | utf8mb4_general_ci |      4 |
// +---------+---------------+--------------------+--------+
// 1 row in set (0.00 sec)
func (conf *DbConf) DSN() (string, error) {
	if conf.Ip == "" {
		conf.Ip = "127.0.0.1"
	}

	if conf.Port == 0 {
		conf.Port = 3306
	}

	if conf.Charset == "" {
		conf.Charset = "utf8mb4"
	}

	// 默认字符序，定义了字符的比较规则
	if conf.Collation == "" {
		conf.Collation = "utf8mb4_general_ci"
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
		return "", err
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

	return mysqlConf.FormatDSN(), nil
}

// InitInstance 建立db连接句柄
// 创建当前数据库db对象，并非连接，在使用的时候才会真正建立db连接
// 为兼容之前的版本，这里新增SetDb创建db对象
func (conf *DbConf) InitInstance() error {
	if conf.hasInit {
		return nil
	}

	// sql日志级别
	if conf.LoggerConfig.LogLevel == 0 {
		conf.LoggerConfig.LogLevel = logger.Info
	}

	// 是否输出sql日志
	// 这里重写了之前的gorm v1版本的日志输出模式
	if conf.ShowSql {
		// 日志对象接口
		var dbLogger logger.Interface
		if conf.Logger == nil {
			dbLogger = logger.Default
		} else {
			dbLogger = logger.New(conf.Logger, conf.LoggerConfig)
		}

		// 设置gorm logger句柄对象
		conf.gormConfig.Logger = dbLogger
	} else {
		conf.gormConfig.Logger = logger.Discard // 默认是不输出sql
	}

	var err error
	if conf.gMysqlConfig.DSN == "" {
		dsn, err := conf.DSN()
		if err != nil {
			log.Println("mysql dsn format error: ", err)
			return err
		}

		conf.gMysqlConfig.DSN = dsn
	}

	// 下面这种方式实例的gorm.DB 很多参数都没法正确设置，不推荐这么实例化
	// conf.dbObj, err = gorm.Open(gMysql.Open(conf.gMysqlConfig.DSN), &gorm.Config{
	// 	Logger: conf.gormConfig.Logger,
	// })

	// 对于golang的官方sql引擎，sql.open并非立即连接db,用的时候才会真正的建立连接
	// 但是gorm.Open在设置完db对象后，还发送了一个Ping操作，判断连接是否连接上去
	// 具体可以看gorm/main.go源码Open方法
	conf.dbObj, err = gorm.Open(gMysql.New(conf.gMysqlConfig), &conf.gormConfig)
	if err != nil {
		log.Println("open mysql connection error: ", err)
		return err
	}

	// 设置连接池
	var sqlDB *sql.DB
	if !conf.hasInit {
		sqlDB, err = conf.SqlDB()
		if err != nil {
			log.Println("get sql db error: ", err)
			return err
		}
	}

	if conf.UsePool {
		sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
		sqlDB.SetMaxOpenConns(conf.MaxOpenConns)
	}

	// 设置连接可以重用的最大存活时间，时间小于数据库的超时时间
	if conf.MaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(conf.MaxLifetime) * time.Second)
	}

	conf.hasInit = true

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

// CloseAllDb 由于gorm db.Close()是关闭当前连接
// 一般建议如下函数放在main/init关闭连接就可以
func CloseAllDb() {
	for name, db := range engineMap {
		sqlDB, err := db.DB()
		if err != nil {
			log.Println("get db instance error: ", err.Error())
			continue
		}

		err = sqlDB.Close()
		if err != nil {
			log.Println("close current db error: ", err)
			continue
		}

		delete(engineMap, name) // 销毁连接句柄标识
	}
}

// CloseDbByName 关闭指定name的db engine
func CloseDbByName(name string) error {
	if db, ok := engineMap[name]; ok {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}

		err = sqlDB.Close()
		if err != nil {
			return err
		}

		delete(engineMap, name) // 销毁连接句柄标识
	}

	return errors.New("current dbObj not exist")
}
