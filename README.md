# starting-grpc
`スターティングgRPC` の勉強リポジトリ

## 本家リポジトリ
https://github.com/gami/grpc_example

## protoc のインストールと実行
```
$ sudo apt install -y protobuf-compiler
$ go get -u github.com/golang/protobuf/protoc-gen-go
$ protoc -Iproto --go_out=plugins=grpc:api --go_opt=module=github.com/hiroygo/starting-grpc/api proto/*.proto
```

## grpc_cli のビルド
```
$ git clone https://github.com/grpc/grpc
$ cd grpc
$ git submodule update --init
$ mkdir -p cmake/build
$ cd cmake/build
$ cmake -DgRPC_BUILD_TESTS=ON ../..
$ make grpc_cli
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
