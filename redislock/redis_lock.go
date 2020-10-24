package redislock

import (
	"github.com/gomodule/redigo/redis"
)

var DefaultExpire = 10 // 加锁的key默认过期时间，单位s

// Lock lock data.
type Lock struct {
	conn   redis.Conn  // redis连接句柄，支持redis pool连接句柄
	expire int         // 设置加锁key的过期时间
	key    string      // 加锁的key
	val    interface{} // 加锁的value
}

// New 实例化redis分布式锁实例对象
func New(conn redis.Conn, key string, val interface{}, expire int) *Lock {
	if expire <= 0 {
		expire = DefaultExpire
	}

	return &Lock{
		key:    key,
		conn:   conn,
		val:    val,
		expire: expire,
	}
}

// delScript lua脚本删除一个key保证原子性，采用lua脚本执行
// 保证原子性（redis是单线程），避免del删除了，其他client获得的lock
var delScript = redis.NewScript(1, `
if redis.call("get", KEYS[1]) == ARGV[1] then
	return redis.call("del", KEYS[1])
else
	return 0
end`)

// Unlock 释放锁采用redis lua脚步执行，成功返回nil
func (lock *Lock) Unlock() error {
	_, err := delScript.Do(lock.conn, lock.key, lock.val)
	return err
}

// TryLock 尝试加锁,如果加锁成功就返回true,nil
// 利用redis setEx nx的原子性实现分布式锁
func (lock *Lock) TryLock() (bool, error) {
	_, err := redis.String(lock.conn.Do("SET", lock.key, lock.val, "EX", lock.expire, "NX"))
	if err != nil {
		return false, err
	}

	return true, nil
}
