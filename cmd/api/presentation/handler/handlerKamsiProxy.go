package api

import (
	"context"
	kamsi "crave/shared/middleware/kamsi/cmd/lib"
	pb "crave/shared/proto/miner"
	"fmt"
	"time"
)

type HandlerKamsiProxy struct {
	kamsi *kamsi.Kamsi
	pb.UnimplementedMinerServer
	hdr IHandler
}

func NewHandlerKamsiProxy(kamsi *kamsi.Kamsi, hdr IHandler) *HandlerKamsiProxy {

	return &HandlerKamsiProxy{kamsi: kamsi, hdr: hdr}
}

func (h *HandlerKamsiProxy) Parse(ctx context.Context, req *pb.ParseRequest) (*pb.ParseResponse, error) {
	start := time.Now()
	res, err := h.hdr.Parse(ctx, req)
	go h.kamsi.Timing("MinerHandler", "Parse", fmt.Sprintf("%s.parseDuration", req.Name), int64(time.Since(start)/time.Millisecond))
	go h.kamsi.Increment("MinerHandler", "Parse", fmt.Sprintf("target.%s.NumbernextTarget", req.Name), int64(len(res.Targets)))
	return res, err
}

// Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
// 	Page          int32                  `protobuf:"varint,2,opt,name=page,proto3" json:"page,omitempty"`
// 	Filter        int64                  `protobuf:"varint,3,opt,name=filter,proto3" json:"filter,omitemp

func (h *HandlerKamsiProxy) Filter(ctx context.Context, req *pb.FilterRequest) (*pb.FilterResponse, error) {
	start := time.Now()
	res, err := h.hdr.Filter(ctx, req)
	go h.kamsi.Timing("MinerHandler", "Filter", fmt.Sprintf("%s.filterDuration", req.Name), int64(time.Since(start)/time.Millisecond))
	go h.kamsi.Increment("MinerHandler", "Filter", fmt.Sprintf("filter.%s", req.Name), int64(len(req.Name)))
	return res, err
}
