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
    
    please look thinkgo code source

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
