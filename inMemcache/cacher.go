package inMemcache

//定义缓存接口
//Set/Get的缓存是存在内存中的，不过期处理
//当前机器重启后，就会丢失缓存
//如果Cacher是基于redis实现的接口，就会具有redis的数据持久化特点
type Cacher interface {
	Get(string) ([]byte, error)
	Set(string, []byte) error
	Delete(string) error
	GetStat() Stat
}

//管理缓存信息的stat
type Stat struct {
	Count     int64 //当前缓冲区key个数
	KeySize   int64 //所有的key大小长度
	ValueSize int64 //所有值的大小长度
}

func (s *Stat) add(k string, v []byte) {
	s.Count += 1
	s.KeySize += int64(len(k))
	s.ValueSize += int64(len(v))
}

func (s *Stat) del(k string, v []byte) {
	s.Count -= 1
	s.KeySize -= int64(len(k))
	s.ValueSize -= int64(len(v))
}
