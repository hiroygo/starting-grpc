package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"github.com/hiroygo/starting-grpc/clientstreaming/api"
	"google.golang.org/grpc"
)

// len(b) が 0 の時は画像本体は送信されない
func uploadImage(c api.ImageUploadServiceClient, filename string, b []byte) (*api.ImageUploadResponse, error) {
	// 画像本体を分割して送信するヘルパー関数
	sendImage := func(stream api.ImageUploadService_UploadClient, chunk int, b []byte) error {
		r := bytes.NewReader(b)
		for {
			// 分割されたデータを読み込む
			lr := io.LimitReader(r, int64(chunk))
			b, err := ioutil.ReadAll(lr)
			if err != nil {
				return fmt.Errorf("ReadAll error: %w", err)
			}
			if len(b) == 0 {
				return nil
			}
			// 分割されたデータを送信する
			image := &api.ImageUploadRequest{
				File: &api.ImageUploadRequest_Data{
					Data: b,
				},
			}
			err = stream.Send(image)
			if err != nil {
				return fmt.Errorf("Send error: %w", err)
			}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream, err := c.Upload(ctx)
	if err != nil {
		return nil, fmt.Errorf("Upload error: %w", err)
	}
	// ファイル情報を送信する
	filemeta := &api.ImageUploadRequest{
		File: &api.ImageUploadRequest_FileMeta_{
			FileMeta: &api.ImageUploadRequest_FileMeta{Filename: filename},
		},
	}
	err = stream.Send(filemeta)
	if err != nil {
		return nil, fmt.Errorf("Send error: %w", err)
	}
	// 画像本体を送信する
	err = sendImage(stream, 1024*10, b)
	if err != nil {
		return nil, fmt.Errorf("sendImage error: %w", err)
	}
	// 返答を受信
	resp, err := stream.CloseAndRecv()
	if err != nil {
		return nil, fmt.Errorf("CloseAndRecv error: %w", err)
	}
	return resp, nil
}

func openFile(path string) (string, []byte, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", nil, fmt.Errorf("ReadFile error: %w", err)
	}
	return filepath.Base(path), b, nil
}

func main() {
	addr := ":50051"
	// grpc.WithInsecure() で TLS ではなく平文で接続する
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatal("Dial error: ", err)
	}
	defer conn.Close()

	filename, image, err := openFile("dog.jpg")
	if err != nil {
		log.Fatal("openFile error: ", err)
	}

	client := api.NewImageUploadServiceClient(conn)
	resp, err := uploadImage(client, filename, image)
	if err != nil {
		log.Fatal("uploadImage error: ", err)
	}
	fmt.Printf("%v\n", resp)
}
