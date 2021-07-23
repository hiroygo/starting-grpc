package handler

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hiroygo/starting-grpc/unary/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type bakedCounter struct {
	// 複数人が同時に焼いても大丈夫なようにしておく
	sync.Mutex
	counts map[api.Pancake_Menu]int
}

func (b *bakedCounter) add(m api.Pancake_Menu, add int) {
	b.Lock()
	defer b.Unlock()
	b.counts[m] = b.counts[m] + add
}

func (b *bakedCounter) report() *api.ReportResponse {
	b.Lock()
	defer b.Unlock()

	counts := []*api.Report_BakeCount{}
	for k, v := range b.counts {
		counts = append(counts, &api.Report_BakeCount{
			Menu: k, Count: int32(v),
		})
	}
	return &api.ReportResponse{Report: &api.Report{BakeCounts: counts}}
}

// pancake.proto の PancakeBakerService に対応する
type BakerHandler struct {
	counter *bakedCounter
}

// インタフェースが実装できていることをコンパイル時に確認する
var _ api.PancakeBakerServiceServer = &BakerHandler{}

func NewBakerHandler() *BakerHandler {
	return &BakerHandler{&bakedCounter{counts: make(map[api.Pancake_Menu]int)}}
}

func (b *BakerHandler) Bake(ctx context.Context, r *api.BakeRequest) (*api.BakeResponse, error) {
	if r.Menu == api.Pancake_UNKNOWN || r.Menu > api.Pancake_SPICY_CURRY {
		return nil, status.Errorf(codes.InvalidArgument, "BakeRequest.Menu error: %v", r.Menu)
	}

	now := time.Now()
	b.counter.add(r.Menu, 1)

	return &api.BakeResponse{
		Pancake: &api.Pancake{
			Menu:           r.Menu,
			ChefName:       "sophia",
			TechnicalScore: rand.Float32(),
			CreateTime: &timestamp.Timestamp{
				Seconds: now.Unix(),
				Nanos:   int32(now.Nanosecond()),
			},
		},
	}, nil
}

func (b *BakerHandler) Report(ctx context.Context, r *api.ReportRequest) (*api.ReportResponse, error) {
	return b.counter.report(), nil
}
