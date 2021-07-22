package main

import (
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/hiroygo/starting-grpc/api"
	"github.com/hiroygo/starting-grpc/server/handler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
)

func main() {
	addr := ":50051"
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Listen error: %s", err)
	}

	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("NewProduction error: %s", err)
	}

	// ReplaceGrpcLogger だけでもロギングは動作する
	// 以下のログだけ出力される
	//{
	//  "level": "info",
	//  "ts": 1626972074.5647643,
	//  "caller": "zap/grpclogger.go:47",
	//  "msg": "[transport]transport: loopyWriter.run returning. connection error: desc = \"transport is closing\"",
	//  "system": "grpc",
	//  "grpc_log": true
	//}
	grpc_zap.ReplaceGrpcLogger(zapLogger)

	sv := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_zap.UnaryServerInterceptor(zapLogger)),
	)
	// gRPC サーバとハンドラを対応させる
	api.RegisterPancakeBakerServiceServer(sv, handler.NewBakerHandler())
	// grpc_cli を使う時に必要
	reflection.Register(sv)
	go func() {
		log.Printf("starting gRPC server: %q\n", addr)
		sv.Serve(l)
	}()

	q := make(chan os.Signal)
	signal.Notify(q, os.Interrupt)
	<-q
	log.Println("stopping gRPC server")
	sv.GracefulStop()
}
