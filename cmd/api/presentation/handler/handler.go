package api

import (
	"context"
	"crave/miner/cmd/api/presentation/controller"
	craveModel "crave/shared/model"
	pb "crave/shared/proto/miner"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Handler struct {
	pb.UnimplementedMinerServer
	ctrl controller.IController
}

func NewHandler(ctrl controller.IController) *Handler {

	return &Handler{ctrl: ctrl}
}

func (h *Handler) Parse(ctx context.Context, req *pb.ParseRequest) (*empty.Empty, error) {
	go h.ctrl.Parse(craveModel.Step(req.Step), craveModel.Page(req.Page), req.Name, craveModel.Filter(req.Filter))
	return &emptypb.Empty{}, status.Error(codes.OK, "")
}
