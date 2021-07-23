## protoc の実行
```
$ protoc -Iproto --go_out=plugins=grpc:api --go_opt=module=github.com/hiroygo/starting-grpc/unary/api proto/*.proto
```

## grpc_cli の実行例
```
$ grpc_cli ls localhost:50051
grpc.reflection.v1alpha.ServerReflection
pancake.baker.PancakeBakerService

$ grpc_cli call localhost:50051 pancake.baker.PancakeBakerService.Bake 'menu: 1'
connecting to localhost:50051
pancake {
  chef_name: "sophia"
  menu: CLASSIC
  technical_score: 0.338417083
  create_time {
    seconds: 1626537673
    nanos: 722801221
  }
}
Rpc succeeded with OK status

$ grpc_cli call localhost:50051 pancake.baker.PancakeBakerService.Report ''
connecting to localhost:50051
report {
  bake_counts {
    menu: CLASSIC
    count: 1
  }
}
Rpc succeeded with OK status
```
