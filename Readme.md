# Thinkgo
    Common library and components for go web and restful api.
# License
    MIT
# Author
    daheige
# About package
    采用govendor机制引入第三包
    .
    ├── cache      go实现cache的存储
    ├── common     公共函数库,包含Time,Lock,Log,redis操作,uuid获取,yaml读取等
    ├── crypto     md5,sha1,sha1File,cbc256,ecb,aes加解密函数等
    ├── GoPool     批量执行task pool池
    ├── http       gin restful定义success,error函数等
    ├── jsoniter   json优化库使用
    ├── Readme.md
    ├── vendor
    └── WatchDog   监控狗,用以监控容易失控的循环或超时
