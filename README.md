# 这是一个流式编程实例

本项目基于grpc、Protocol Buffers，需要安装好protoc工具

protoc version：3.21.12

还需要安装最新的protoc-gen-go
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

环境变量设置好go bin和go path。

详情见 https://grpc.io/docs/languages/go/quickstart/

# 试用
先运行go_web_srv
再运行go_web_client