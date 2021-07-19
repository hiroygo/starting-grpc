package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/hiroygo/starting-grpc/api"
	"github.com/hiroygo/starting-grpc/server/handler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := 50051
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Listen error: %s", err)
	}

	sv := grpc.NewServer()
	// gRPC サーバとハンドラを対応させる
	api.RegisterPancakeBakerServiceServer(sv, handler.NewBakerHandler())
	// grpc_cli を使う時に必要
	reflection.Register(sv)
	go func() {
		log.Printf("starting gRPC server: %d\n", port)
		sv.Serve(l)
	}()

	q := make(chan os.Signal)
	signal.Notify(q, os.Interrupt)
	<-q
	log.Println("stopping gRPC server")
	sv.GracefulStop()
}
