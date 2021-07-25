## プログラム概要
gRPC を使用せずに protocol-buffers のシリアライズ/デシリアライズだけ試す

## protoc の実行
`-go_out=plugins=grpc:` オプションを付けると gRPC 用のコードが生成される

```
$ protoc --go_out=. ./dog.proto --go_opt=module=github.com/hiroygo/starting-grpc/protocol-buffers/main
```

## シリアライズされたデータを見てみる
```
$ ./protocol-buffers
name:"sophia" age:7
$ od -tx dog
0000000 6f73060a 61696870 00000710
0000012
```
