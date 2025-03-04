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

	h.kamsi.Timing("MinerHandler", "Parse", fmt.Sprintf("target.%s.parse", req.Name), int64(time.Since(start)/time.Millisecond))
	h.kamsi.IntGauge("MinerHandler", "Parse", fmt.Sprintf("target.%s.parse.nextTarget", req.Name), func(err error) int64 {
		if err != nil {
			return -1
		}
		return int64(len(res.Targets))
	}(err))
	return res, err
}

// Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
// 	Page          int32                  `protobuf:"varint,2,opt,name=page,proto3" json:"page,omitempty"`
// 	Filter        int64                  `protobuf:"varint,3,opt,name=filter,proto3" json:"filter,omitemp

func (h *HandlerKamsiProxy) Filter(ctx context.Context, req *pb.FilterRequest) (*pb.FilterResponse, error) {
	start := time.Now()
	res, err := h.hdr.Filter(ctx, req)
	h.kamsi.Timing("MinerHandler", "Filter", fmt.Sprintf("target.%s.filter", req.Name), int64(time.Since(start)/time.Millisecond))
	h.kamsi.Increment("MinerHandler", "Filter", fmt.Sprintf("target.%s.filter.response", req.Name), 1)
	return res, err
}
