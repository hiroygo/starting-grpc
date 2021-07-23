## protoc の実行
```
$ protoc -Iproto --go_out=plugins=grpc:api --go_opt=module=github.com/hiroygo/starting-grpc/clientstreaming/api proto/*.proto
```

