package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/hiroygo/starting-grpc/api"
	"github.com/hiroygo/starting-grpc/server/handler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
)

func newServer() (*grpc.Server, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("NewProduction error: %w", err)
	}

	// 認証用処理
	auth := func(ctx context.Context) (context.Context, error) {
		// MD = metadata
		token, err := grpc_auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			return nil, fmt.Errorf("AuthFromMD error: %w", err)
		}
		if token != "secret" {
			return nil, grpc.Errorf(codes.Unauthenticated, "invalid bearer token")
		}
		return ctx, nil
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
	grpc_zap.ReplaceGrpcLogger(logger)
	sv := grpc.NewServer(
		grpc.UnaryInterceptor(
			// ChainUnaryServer(one, two, three) は one, two, three の順で実行されていく
			grpc_middleware.ChainUnaryServer(
				grpc_zap.UnaryServerInterceptor(logger),
				grpc_auth.UnaryServerInterceptor(auth),
			),
		),
	)
	return sv, nil
}

func main() {
	addr := ":50051"
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Listen error: %s", err)
	}

	sv, err := newServer()
	if err != nil {
		log.Fatalf("newServer error: %s", err)
	}

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
