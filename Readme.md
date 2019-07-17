# Thinkgo
    Public libraries and components for glang development.
# Package manage
    golang1.11以下版本，采用govendor机制引入第三包或直接使用vendor
    golang1.11+版本，可采用go mod机制管理包

# About package
    .
    ├── bitSet          bitSet位图实现，将int数放到内存中进行curd操作
    ├── cache           基于key/val内存缓存设计，支持过期时间设置
    ├── common          公共函数库logger每天流动日志,Time,Lock操作,uuid生成等
    ├── crypto          md5,sha1,sha1File,cbc256,ecb,aes加解密函数等
    ├── gQueue          通过指定goroutine个数,实现task queue执行器
    ├── httpRequest     http request请求封装，支持get,post,put,patch,delete等等
    ├── inMemcache      通过接口的形式实现内存cache实现kv存储
    ├── logger          基于uber zap框架封装而成的高性能logger日志库
    ├── mysql           gorm/mysql封装,gorm1.9.9+版本
    ├── dao             xorm/mysql封装,支持读写分离连接对象设置
    ├── mytest          单元测试用例
    ├── rbmq            rbmq连接封装
    ├── redisCache      redisgo操作库封装
    ├── runner          runner按照顺序，执行任务操作，可作为cron作业或定时任务
    ├── work            利用无缓冲chan创建goroutine池来控制一组task的执行
    ├── workPool        workerPool工作池，实现百万级的并发,一般用于持续不断的大规模作业
    ├── xerrors         error错误处理拓展包，支持错误堆栈信息
    └── yamlConf        yaml配置文件读取

# Use help
    设置golang proxy
    vim ~/.bashrc添加如下内容：
    export GOPROXY=https://goproxy.io
    或者使用 export GOPROXY=https://athens.azurefd.net
    让bashrc生效
    source ~/.bashrc
    
    如果是采用govendor管理包请按照如下方式进行：
        1. 下载thinkgo包
            cd $GOPATH/src
            git clone https://github.com/daheige/thinkgo.git
        2. 安装govendor go第三方包管理工具
            go get -u github.com/kardianos/govendor
        3. 切换到对应的目录进行 go install编译包
    如果采用go mod (golang1.11版本+) 不需要将该包放在$GOPATH/src，只需要在使用的项目中引入就可以。
        请直接执行go mod tidy # 下载依赖包
# License
    MIT
