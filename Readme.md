# Thinkgo package

    Public libraries and components for glang development.

    I like the language of php. I have been using php development experience for 5 years.
    It has inspired me a lot. I quickly converted to golang development in 2 years.
    I am very glad to be exposed to this language.
    These functions and packages are used extensively in development,
    so they are packaged as components or libraries for development.

# Package management

    golang1.11+版本，可采用go mod机制管理包,需设置goproxy
    golang1.11以下版本，采用govendor机制引入第三包或直接使用vendor

# About package

    .
    ├── bitset                      bitSet位图实现，将int数放到内存中进行curd操作
    ├── common                      公共函数库，包含每天流动日志,Time,Lock操作,uuid生成等
    ├── crypto                      md5,sha1,sha1File,cbc256,ecb,aes加解密函数等
    ├── gnsq                        go nsq消费队列封装,主要是pub/sub模式，其他类型也支持句柄调用
    ├── goredis                     基于go-redis/redis封装的redis客户端操作，支持redis cluster集群模式
    ├── gqueue                      通过指定goroutine个数,实现task queue执行器
    ├── gresty                      request请求封装，支持get,post,put,patch,delete等
    ├── gxorm                       基于xorm/mysql封装,支持读写分离连接对象设置
    ├── jsontime                    fix time.Time datetime格式的json encode/decode bug
    ├── logger                      基于uber zap框架封装而成的高性能logger日志库
    ├── monitor                     用于对go程序做prometheus/metrics pprof性能监控，包含内存，cpu,请求数等
    ├── mysql                       gorm/mysql封装,主要基于gorm1.9.10+版本
    ├── mytest                      thinkgo单元测试用例
    ├── rediscache                  redisgo操作库封装
    ├── redislock                   redis+lua脚步实现redis分布式锁TryLock,Unlock
    ├── runner                      runner按照顺序，执行任务操作，可作为cron作业或定时任务
    ├── sem                         指定数量的空结构体缓存通道，实现信息号实现互斥锁
    ├── work                        利用无缓冲chan创建goroutine池来控制一组task的执行
    ├── workpool                    workerPool工作池，实现百万级的并发,一般用于持续不断的大规模作业
    ├── xerrors                     error错误处理拓展包，支持错误堆栈信息
    ├── xsort                       基于标准包sort封装的一些sort操作方法
    ├── yamlconf                    yaml配置文件读取，支持int,int64,float64,string,struct等类型读取
    ├── common/str_convert.go       字符串，数字相互转换的一些辅助函数，参考php函数实现
    ├── common/file.go              文件/目录相关的函数，参考php函数实现
    └── common/num.go               数字相关的函数，参考php函数实现

# usage

    go version >= 1.13
    设置goproxy代理
    vim ~/.bashrc添加如下内容:
    export GOPROXY=https://goproxy.io,direct
    或者
    export GOPROXY=https://goproxy.cn,direct
    或者
    export GOPROXY=https://mirrors.aliyun.com/goproxy/,direct

    让bashrc生效
    source ~/.bashrc

    go version < 1.13
    设置golang proxy
    vim ~/.bashrc添加如下内容：
    export GOPROXY=https://goproxy.io
    或者使用 export GOPROXY=https://athens.azurefd.net
    或者使用 export GOPROXY=https://mirrors.aliyun.com/goproxy/ #推荐该goproxy
    让bashrc生效
    source ~/.bashrc

    go version < 1.11
    如果是采用govendor管理包请按照如下方式进行：
        1. 下载thinkgo包
            cd $GOPATH/src
            git clone https://github.com/daheige/thinkgo.git
        2. 安装govendor go第三方包管理工具
            go get -u github.com/kardianos/govendor
        3. 切换到对应的目录进行 go install编译包

# Test unit

    测试mytest
    $ go test -v
    997: b75567dc6f88412d55576e4b09127d3f
    998: c3923160f2304849734c0907083f7f65
    999: 8b7a6dce56d346b567c65b3493285831
    --- PASS: TestUuid (0.05s)
        uuid_test.go:13: 测试uuid
    PASS
    ok      github.com/daheige/thinkgo/mytest       15.841s

    $ cd common
    $ go test -v
    2019/10/28 22:32:01 current rnd uuid a3e96dae-ca2a-d029-76c9-279b1fff1234
    2019/10/28 22:32:01 current rnd uuid 4e136db3-56a8-fa67-93d7-f11f6cfd57ae
    2019/10/28 22:32:01 current rnd uuid 30b83e42-2040-7ab3-9089-05d0d558bbcc
    2019/10/28 22:32:01 current rnd uuid 16bc0ad1-4b17-27ee-2a2d-7b08f175295b
    2019/10/28 22:32:01 current rnd uuid 979aefef-9db9-baad-920a-1742d24c2166
    --- PASS: TestRndUuid (35.71s)
    PASS
    ok  	github.com/daheige/thinkgo/common	91.290s

# License

    MIT
