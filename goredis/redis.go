package goredis

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis"
)

// a redis client list
var RedisClientList = map[string]*redis.Client{}

var HashDefaultExpire int64 = 300 // 默认过期时间300s

// redis client config
type RedisClientConf struct {
	// host:port address.
	Address string

	// Optional password. Must match the password specified in the
	// requirepass server configuration option.
	Password string

	// Database to be selected after connecting to the server.
	DB int

	// Maximum number of retries before giving up.
	// Default is to not retry failed commands.
	MaxRetries int

	// Dial timeout for establishing new connections.
	// Default is 5 seconds.
	DialTimeout time.Duration

	// Timeout for socket reads. If reached, commands will fail
	// with a timeout instead of blocking. Use value -1 for no timeout and 0 for default.
	// Default is 3 seconds.
	ReadTimeout time.Duration

	// Timeout for socket writes. If reached, commands will fail
	// with a timeout instead of blocking.
	// Default is ReadTimeout.
	WriteTimeout time.Duration

	// Maximum number of socket connections.
	// Default is 10 connections per every CPU as reported by runtime.NumCPU.
	PoolSize int

	// Amount of time client waits for connection if all connections
	// are busy before returning an error.
	// Default is ReadTimeout + 1 second.
	PoolTimeout time.Duration

	// Minimum number of idle connections which is useful when establishing
	// new connection is slow.
	MinIdleConns int

	// Amount of time after which client closes idle connections.
	// Should be less than server's timeout.
	// Default is 5 minutes. -1 disables idle timeout check.
	IdleTimeout time.Duration

	// Connection age at which client retires (closes) the connection.
	// Default is to not close aged connections.
	MaxConnAge time.Duration
}

// redis cluster config
type RedisClusterConf struct {
	// A seed list of host:port addresses of cluster nodes.
	AddressNodes []string

	Password string

	// Maximum number of retries before giving up.
	// Default is to not retry failed commands.
	MaxRetries int

	DialTimeout  time.Duration // Default is 5 seconds.
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	// PoolSize applies per cluster node and not for the whole cluster.
	PoolSize int

	// Amount of time client waits for connection if all connections
	// are busy before returning an error.
	// Default is ReadTimeout + 1 second.
	PoolTimeout time.Duration

	// Minimum number of idle connections which is useful when establishing
	// new connection is slow.
	MinIdleConns int

	// Amount of time after which client closes idle connections.
	// Should be less than server's timeout.
	// Default is 5 minutes. -1 disables idle timeout check.
	IdleTimeout time.Duration

	// Connection age at which client retires (closes) the connection.
	// Default is to not close aged connections.
	MaxConnAge time.Duration
}

// GetClient return redis client
func (conf *RedisClientConf) GetClient() *redis.Client {
	if conf.MaxConnAge == 0 {
		conf.MaxConnAge = 30 * 60 * time.Second
	}

	opt := &redis.Options{
		Addr:         conf.Address,
		Password:     conf.Password,
		DB:           conf.DB, // use default DB
		MaxRetries:   conf.MaxRetries,
		DialTimeout:  conf.DialTimeout,  // Default is 5 seconds.
		ReadTimeout:  conf.ReadTimeout,  // Default is 3 seconds.
		WriteTimeout: conf.WriteTimeout, // Default is ReadTimeout.
		PoolSize:     conf.PoolSize,
		PoolTimeout:  conf.PoolTimeout,
		MinIdleConns: conf.MinIdleConns,
		IdleTimeout:  conf.IdleTimeout,
		MaxConnAge:   conf.MaxConnAge,
	}

	return redis.NewClient(opt)
}

// SetClientName set a redis client to clientList
func (conf *RedisClientConf) SetClientName(name string) {
	RedisClientList[name] = conf.GetClient()
}

// GetRedisClient get redis client from RedisClientList
func GetRedisClient(name string) (*redis.Client, error) {
	if _, ok := RedisClientList[name]; ok {
		return RedisClientList[name], nil
	}

	return nil, errors.New("current client " + name + " not exist")
}

// SetJson 设置任意类型到redis中，以json格式保存
func SetJson(client *redis.Client, key string, d interface{}, expire int64) error {
	b, err := json.Marshal(d)
	if err != nil {
		return err
	}

	if expire > 0 {
		_, err = client.Do("setEx", key, expire, string(b)).Result()
	} else {
		_, err = client.Do("set", key, string(b)).Result()
	}

	return err
}

// GetJson 从redis中获取指定的key对应的val解析到data中
// data必须是提前定义好的类型，且为指针类型
func GetJson(client *redis.Client, key string, data interface{}) error {
	str, err := client.Do("get", key).String()
	if err != nil {
		return err
	}

	if str == "" {
		return errors.New("redis data is empty")
	}

	err = json.Unmarshal([]byte(str), data)
	if err != nil {
		return err
	}

	return nil
}

// GetCluster return redis cluster client
func (conf *RedisClusterConf) GetCluster() *redis.ClusterClient {
	if conf.MaxConnAge == 0 {
		conf.MaxConnAge = 30 * 60 * time.Second
	}

	cluster := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        conf.AddressNodes,
		Password:     conf.Password,
		PoolSize:     conf.PoolSize,
		MaxRetries:   conf.MaxRetries,
		DialTimeout:  conf.DialTimeout,
		ReadTimeout:  conf.ReadTimeout,
		WriteTimeout: conf.WriteTimeout,
		PoolTimeout:  conf.PoolTimeout,
		MinIdleConns: conf.MinIdleConns,
		IdleTimeout:  conf.IdleTimeout,
		MaxConnAge:   conf.MaxConnAge,
	})

	return cluster
}
