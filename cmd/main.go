package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hiroygo/starting-grpc/api"
	"google.golang.org/grpc"
)

func bake(c api.PancakeBakerServiceClient, m api.Pancake_Menu) (*api.BakeResponse, error) {
	req := &api.BakeRequest{Menu: m}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	resp, err := c.Bake(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("Bake error: %w", err)
	}

	return resp, nil
}

func report(c api.PancakeBakerServiceClient) (*api.ReportResponse, error) {
	req := &api.ReportRequest{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	resp, err := c.Report(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("Report error: %w", err)
	}

	return resp, nil
}

func main() {
	addr := ":50051"
	// grpc.WithInsecure() で TLS ではなく平文で接続する
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("Dial error: ", err)
	}
	defer conn.Close()

	client := api.NewPancakeBakerServiceClient(conn)
	_, err = bake(client, api.Pancake_CLASSIC)
	if err != nil {
		log.Fatal("bake error: ", err)
	}
	report, err := report(client)
	if err != nil {
		log.Fatal("report error: ", err)
	}
	fmt.Printf("%v\n", report.Report)
}
