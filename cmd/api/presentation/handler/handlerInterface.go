package api

import (
	pb "crave/shared/proto/miner"
)

type IHandler interface {
	pb.MinerServer
}
