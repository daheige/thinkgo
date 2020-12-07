// Package mysql gMysql option for gorm v2 config.
package mysql

import (
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Option DbConf 功能函数模式
type Option func(conf *DbConf)

// WithDriverName 设置db driver name.
func WithDriverName(name string) Option {
	return func(conf *DbConf) {
		conf.gMysqlConfig.DriverName = name
	}
}

// WithDsn 设置dsn
func WithDsn(dsn string) Option {
	return func(conf *DbConf) {
		conf.gMysqlConfig.DSN = dsn
	}
}

// WithGormConnPool 设置gorm conn pool.
func WithGormConnPool(conn gorm.ConnPool) Option {
	return func(conf *DbConf) {
		conf.gMysqlConfig.Conn = conn
	}
}

// WithSkipInitializeWithVersion 根据当前 MySQL 版本自动配置
func WithSkipInitializeWithVersion(b bool) Option {
	return func(conf *DbConf) {
		conf.gMysqlConfig.SkipInitializeWithVersion = b
	}
}

// WithDefaultStringSize 设置字符串默认长度
func WithDefaultStringSize(size uint) Option {
	return func(conf *DbConf) {
		conf.gMysqlConfig.DefaultStringSize = size
	}
}

// WithDisableDatetimePrecision 是否禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
func WithDisableDatetimePrecision(b bool) Option {
	return func(conf *DbConf) {
		conf.gMysqlConfig.DisableDatetimePrecision = b
	}
}

// WithDontSupportRenameIndex 是否重命名索引时采用删除并新建的方式
// MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
func WithDontSupportRenameIndex(b bool) Option {
	return func(conf *DbConf) {
		conf.gMysqlConfig.DontSupportRenameIndex = b
	}
}

// WithDontSupportRenameColumn
// 是否用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
func WithDontSupportRenameColumn(b bool) Option {
	return func(conf *DbConf) {
		conf.gMysqlConfig.DontSupportRenameColumn = b
	}
}

// WithGormConfig 设置gorm.Config
func WithGormConfig(config gorm.Config) Option {
	return func(conf *DbConf) {
		conf.gormConfig = config
	}
}

// WithLogger 设置gorm logger实例对象
func WithLogger(gLogger logger.Writer) Option {
	return func(conf *DbConf) {
		conf.Logger = gLogger
	}
}

// WithLogLevel 设置sql logger level
func WithLogLevel(logLevel logger.LogLevel) Option {
	return func(conf *DbConf) {
		conf.LoggerConfig.LogLevel = logLevel
	}
}
