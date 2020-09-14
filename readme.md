# Thinkgo package

    Public libraries and components for glang development.

    I like the language of php. I have been using php development experience for 6 years.
    It has inspired me a lot. I quickly converted to golang development in 2 years.
    I am very glad to be exposed to this language.
    These functions and packages are used extensively in development,
    so they are packaged as components or libraries for development.

# Package management

    golang1.11+版本，可采用go mod机制管理包,需设置goproxy
    golang1.11以下版本，采用govendor机制引入第三包或直接使用vendor

# About package
    
    .
    ├── bitset              bitSet位图实现
    ├── chanlock            chan实现trylock乐观锁
    ├── crypto              常见的md5,sha1,sha1file,aes/des,ecb,openssl_encrypt实现
    ├── def                 为兼容php其他语言而定义的空数组，空对象
    ├── gfile               file文件操作的一些辅助函数
    ├── glog                基于mutex乐观锁实现的每天流动式日志，将日志内容直接落地到文件中
    ├── gnsq                go-nsq基本操作封装
    ├── gnum                num Round,Floor,Ceil等函数实现
    ├── goredis             基于go-redis/redis封装的redis客户端使用函数（支持cluster集群）
    ├── gpprof              pprof性能分析监控封装
    ├── gqueue              通过指定goroutine个数,实现task queue执行器
    ├── grecover            golang panic/recover捕获堆栈信息实现
    ├── gresty              go http client support get,post,delete,patch,put,head,file method
    ├── gtask               golang task在独立协程中调度实现
    ├── gtime               time相关的一些辅助函数
    ├── gutils              字符串相关的一些辅助函数，比如Uuid,HTMLSpecialchars,Uniqid等php函数实现
    ├── gxorm               golang xorm客户端简单封装，方便使用
    ├── jsontime            fix gorm/xorm time.Time json encode/decode bug
    ├── logger              基于zap日志库进行一些必要的优化的日志库
    ├── monitor             基于prometheus二次开发、封装的一些函数，主要用于http/job/grpc服务性能监控
    ├── mutexlock           基于sync.Mutex基础上拓展的乐观锁
    ├── mysql               基于go gorm库封装而成的mysql客户端的一些辅助函数
    ├── mytest              thinkgo 一些单元测试
    ├── gredigo             基于redigo封装而成的go redis辅助函数，方便快速接入redis操作
    ├── redislock           基于redigo实现的redis+lua分布式锁实现
    ├── runner              runner用于按照顺序，执行程序任务操作，可作为cron作业或定时任务
    ├── sem                 指定数量的空结构体缓存通道，实现信息号实现互斥锁
    ├── setting             通过viper+fsnotify实现配置文件读取，支持配置热更新
    ├── strlist             string list实现
    ├── work                利用无缓冲chan创建goroutine池来控制一组task的执行
    ├── workpool            workpool工作池实现，对于百万级并发的一些场景特别适用
    ├── xerrors             自定义错误类型，一般用在api/微服务等业务逻辑中，处理错误
    ├── xsort               基于sort标准库封装的sort操作函数
    └── yamlconf            基于yaml+reflect实现yaml文件的读取，一般用在web/job/rpc应用中

# Upgrade log
    
    2020.09.14
        1) add viper read config.
            
    2020.09.11
        1) xorm升级到v1.0.5
        2) gorm升级到v1.9.16
            
    2020.08.30
        1）对xorm从v0.8.2升级到v1.0.3，支持mysql5.6-mysql8.0+版本
        2）对gxorm/gorm mysql sql日志输出采用接口方式设计
        3）废弃gxorm/gorm mysql SqlCmd参数，改为ShowSql
        4）删除gxorm ShowExecTime参数配置
        如果需要使用原来的版本，请使用thinkgo v1.10.x版本

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
    
# License

    MIT
