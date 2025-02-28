package api

import (
	"context"
	"crave/miner/cmd/api/presentation/controller"
	craveModel "crave/shared/model"
	pb "crave/shared/proto/miner"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	pb.UnimplementedMinerServer
	ctrl controller.IController
}

func NewHandler(ctrl controller.IController) *Handler {

	return &Handler{ctrl: ctrl}
}

func (h *Handler) Parse(ctx context.Context, req *pb.ParseRequest) (*pb.ParseResponse, error) {
	names, err := h.ctrl.Parse(craveModel.Step(req.Step), craveModel.Page(req.Page), req.Name)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to parse targets")
	}
	return &pb.ParseResponse{Targets: names}, nil
}

// Name          string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
// 	Page          int32                  `protobuf:"varint,2,opt,name=page,proto3" json:"page,omitempty"`
// 	Filter        int64                  `protobuf:"varint,3,opt,name=filter,proto3" json:"filter,omitemp

func (h *Handler) Filter(ctx context.Context, req *pb.FilterRequest) (*pb.FilterResponse, error) {
	filterBy, err := h.ctrl.Filter(req.Name, craveModel.Page(req.Page), craveModel.Filter(req.Filter))
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to parse targets")
	}
	return &pb.FilterResponse{FilteredBy: filterBy}, err
}
