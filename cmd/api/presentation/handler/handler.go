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
