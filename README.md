# starting-grpc
`スターティングgRPC` の勉強リポジトリ

## 本家リポジトリ
https://github.com/gami/grpc_example

## protoc のインストール
```
$ sudo apt install -y protobuf-compiler
$ go get -u github.com/golang/protobuf/protoc-gen-go
```

## HTML の生成
```
$ go get -u github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc
$ protoc -Iproto --doc_out=. --doc_opt=html,index.html proto/*.proto
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

