# Thinkgo
    Common library and components for go web and restful api.
# govendor manage go package
    采用govendor机制引入第三包
# About package
    .
    ├── bitSet      bitSet位图实现，将int数放到内存中进行curd操作
    ├── cache       go实现cache的存储
    ├── common      公共函数库,包含Time,Lock,Log,redis操作,uuid生成,yaml读取等
    ├── crypto      md5,sha1,sha1File,cbc256,ecb,aes加解密函数等
    ├── GoPool      批量执行task pool池
    ├── gorm        gorm/mysql封装
    ├── gqueue      通过指定goroutine个数,实现task queue执行器
    ├── http        gin restful定义success,error函数等
    ├── jsoniter    json优化库使用     
    ├── rbmq        rbmq连接封装
    ├── runner      runner按照顺序，执行任务操作，可作为cron作业或定时任务  
    ├── WatchDog    监控狗,用以监控容易失控的循环或超时
    └── work        利用无缓冲chan创建goroutine池来控制一组task的执行
# use help
    1. 下载thinkgo包
        cd $GOPATH
        git clone https://github.com/daheige/thinkgo.git
    2. 安装govendor go第三方包管理工具
        go get -u github.com/kardianos/govendor
    3. 切换到对应的目录进行 go install编译包
# License
    MIT
