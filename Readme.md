# Thinkgo
    Common library and components for go web and restful api.
# Author
    daheige
# vendor第三包
    采用govendor机制引入第三包
# About package
    .
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

# License
    MIT
