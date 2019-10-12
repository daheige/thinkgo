# gResty package

        go http client library,base on go-resty.
        go-resty: https://github.com/go-resty/resty

# usage

        go version 1.11+

        import "github.com/daheige/thinkgo/gResty"

        s := &gResty.Service{
            BaseUri: "http://localhost:1338/",
            Timeout: 2 * time.Second,
        }

        opt := &gResty.ReqOpt{
            Data: map[string]interface{}{
                "id": "1234",
            },
            RetryCount:2,
        }

        res := s.Do("post", "v1/data", opt)
        log.Println("err: ", res.Err)
        log.Println("body:", res.Text())

        For other usage, please see the method in the gResty source package.
