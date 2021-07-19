package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hiroygo/starting-grpc/api"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal("Dial error: ", err)
	}
	defer conn.Close()

	request := &api.BakeRequest{Menu: api.Pancake_CLASSIC}
	client := api.NewPancakeBakerServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	response, err := client.Bake(ctx, request)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("%v\n", response.Pancake)
}
