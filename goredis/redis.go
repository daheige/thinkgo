package goredis

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis"
)

// a redis client list
var RedisClientList = map[string]*redis.Client{}

var HashDefaultExpire int64 = 300 //默认过期时间300s

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
}

// GetClient return redis client
func (conf *RedisClientConf) GetClient() *redis.Client {
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

type RedisData struct {
	Expire int64
	Data   interface{}
}

// SetJson 设置任意类型到redis中，以json格式保存
func SetJson(client *redis.Client, key string, d RedisData) error {
	b, err := json.Marshal(d)
	if err != nil {
		return err
	}

	if d.Expire > 0 {
		_, err = client.Do("setEx", key, d.Expire, string(b)).Result()
	} else {
		_, err = client.Do("set", key, string(b)).Result()
	}

	return err
}

// GetJson 从redis中获取指定的key对应的val
func GetJson(client *redis.Client, key string) (interface{}, error) {
	str, err := client.Do("get", key).String()
	if err != nil {
		return nil, err
	}

	if str == "" {
		return nil, errors.New("redis data is empty")
	}

	res := &RedisData{}
	err = json.Unmarshal([]byte(str), res)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

// SetExHash 拓展redis hash里面的key没有过期时间的问题
func SetExHash(client *redis.Client, key string, innerKey string, d RedisData) error {
	if d.Expire <= 0 {
		d.Expire = time.Now().Unix() + HashDefaultExpire
	}

	b, err := json.Marshal(d)
	if err != nil {
		return err
	}

	_, err = client.Do("hSet", key, innerKey, string(b)).Result()

	return err
}

// GetExHash 获取redis hash内部key对应的val
func GetExHash(client *redis.Client, key string, innerKey string) (interface{}, error) {
	str, err := client.Do("hget", key, innerKey).String()
	if err != nil {
		return nil, err
	}

	if str == "" {
		return nil, errors.New("redis data is empty")
	}

	res := &RedisData{}
	err = json.Unmarshal([]byte(str), res)
	if err != nil {
		return nil, err
	}

	//当数据没有过期，就直接返回,否则就删除
	if res.Expire >= time.Now().Unix() {
		return res.Data, nil
	}

	client.Do("hDel", key, innerKey)

	return nil, errors.New("hash data has expired")
}

// GetCluster return redis cluster client
func (conf *RedisClusterConf) GetCluster() *redis.ClusterClient {
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
	})

	return cluster
}
