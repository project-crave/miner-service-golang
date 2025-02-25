package externalApi

import (
	"context"
	"crave/shared/common/client"
	craveModel "crave/shared/model"
	pb "crave/shared/proto/hub"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HubClient struct {
	client     client.IClient
	grpcClient pb.HubClient
}

func NewHubClient(baseURL string, grpcClient pb.HubClient) *HubClient {
	return &HubClient{
		client:     client.NewGenericClient(baseURL),
		grpcClient: grpcClient,
	}
}

func (c *HubClient) ParseResult(name string, targets []string, step craveModel.Step) error {
	ctx := context.Background()
	opts := []grpc.CallOption{
		grpc.UseCompressor("gzip-level-9"),
	}
	_, err := c.grpcClient.ParseResult(ctx, &pb.ParseResultRequest{
		Name:    name,
		Targets: targets,
		Step:    int32(step),
	}, opts...)
	if err != nil {
		return err
	}
	if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.OK:
		default:

		}
	}
	return nil
}
