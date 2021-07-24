package main

import (
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/hiroygo/starting-grpc/clientstreaming/api"
	"github.com/hiroygo/starting-grpc/clientstreaming/server/handler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func newServer() (*grpc.Server, error) {
	return grpc.NewServer(), nil
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
	api.RegisterImageUploadServiceServer(sv, handler.NewImageUploadHandler())
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
