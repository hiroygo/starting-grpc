package handler

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

    "github.com/google/uuid"
	"github.com/hiroygo/starting-grpc/clientstreaming/api"
)

type ImageUploadHandler struct {
	sync.Mutex
	files map[string][]byte
}

func NewImageUploadHandler() *ImageUploadHandler {
	return &ImageUploadHandler{files: make(map[string][]byte)}
}

func (i *ImageUploadHandler) Upload(stream api.ImageUploadService_UploadServer) error {
	recvFilename := func() (string, error) {
		req, err := stream.Recv()
		if err != nil {
			return "", fmt.Errorf("Recv error: %w", err)
		}
		meta := req.GetFileMeta()
		if meta == nil {
			return "", errors.New("GetFileMeta returned nil")
		}
		return meta.Filename, nil
	}
	recvImage := func() ([]byte, error) {
		b := &bytes.Buffer{}
		// 分割されたバイナリをループして受け取る
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				return b.Bytes(), nil
			}
			if err != nil {
				return nil, fmt.Errorf("Recv error: %w", err)
			}

			chunk := req.GetData()
			if chunk == nil {
				return nil, errors.New("GetData returned nil")
			}
			_, err = b.Write(chunk)
			if err != nil {
				return nil, fmt.Errorf("Write error: %w", err)
			}
		}
	}
	sendResponse := func(filename string, image []byte) error {
		id, err := uuid.NewRandom()
		if err != nil {
			fmt.Errorf("NewRandom error: %w", err)
		}
		resp := &api.ImageUploadResponse{
			Uuid:        id.String(),
			Size:        int32(len(image)),
			Filename:    filename,
			ContentType: http.DetectContentType(image),
		}

		err = stream.SendAndClose(resp)
		if err != nil {
			return fmt.Errorf("SendAndClose error: %w", err)
		}
		return nil
	}

	// 通信の初めにメタ情報を受け取る
	filename, err := recvFilename()
	if err != nil {
		return fmt.Errorf("recvFilename error: %w", err)
	}
	// 画像を受け取る
	image, err := recvImage()
	if err != nil {
		return fmt.Errorf("recvImage error: %w", err)
	}
	i.files[filename] = image

	// 返信
	err = sendResponse(filename, image)
	if err != nil {
		return fmt.Errorf("sendResponse error: %w", err)
	}
	return nil
}
