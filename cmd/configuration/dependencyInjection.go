package configuration

import (
	filterBusiness "crave/miner/cmd/api/domain/business/filter"
	pageBusiness "crave/miner/cmd/api/domain/business/page"
	service "crave/miner/cmd/api/domain/service"
	"crave/miner/cmd/api/infrastructure/externalApi"
	pageApi "crave/miner/cmd/api/infrastructure/externalApi/page"
	"crave/miner/cmd/api/infrastructure/repository"
	controller "crave/miner/cmd/api/presentation/controller"
	handler "crave/miner/cmd/api/presentation/handler"
	craveDatabase "crave/shared/database"
	kamsi "crave/shared/middleware/kamsi/cmd/configuration"

	craveModel "crave/shared/model"
	craveGzip "crave/shared/proto/gzip"
	hubPb "crave/shared/proto/hub"
	minerPb "crave/shared/proto/miner"
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Container struct {
	Router          *gin.Engine
	Variable        *Variable
	KamsiCtnr       *kamsi.Container
	MinerController controller.IController
	MinerHandler    handler.IHandler
	MinerService    service.IService
	HubGrpcClient   hubPb.HubClient
	FilterStrategy  *filterBusiness.FilterStrategy
	PageStrategy    *pageBusiness.PageStrategy
	MinerRepository repository.IRepository
	HubClient       externalApi.IHubClient
	Neo4j           *craveDatabase.Neo4jWrapper
}

func (ctnr *Container) SetRouter(router any) {
	if r, ok := router.(*gin.Engine); ok {
		ctnr.Router = r
		return
	}
	panic("Provided router is not a *gin.Engine")
}

func (ctnr *Container) DefineDatabase() error {
	ctnr.Neo4j = craveDatabase.ConnectNeo4jDatabase(ctnr.Variable.Neo4jMiner)
	return nil
}

func (ctnr *Container) DefineGrpc() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ctnr.Variable.GrpcApi.Ip, ctnr.Variable.GrpcApi.Port))
	if err != nil {
		return fmt.Errorf("failed to listen : %d, %w", ctnr.Variable.GrpcApi.Port, err)
	}
	server := grpc.NewServer()
	minerPb.RegisterMinerServer(server, ctnr.MinerHandler)
	if servErr := server.Serve(lis); servErr != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}
	return nil
}

func (ctnr *Container) InitVariable() error {
	ctnr.Variable = NewVariable()
	return nil
}

func (ctnr *Container) DefineRoute() error {
	// minerGroup := ctnr.Router.Group("/miner")
	// {
	// 	minerGroup.POST("/parse", ctnr.MinerHandler.Parse)
	// }
	return nil
}

func (ctnr *Container) GetHttpHandler() http.Handler {
	return nil
}

func (ctnr *Container) InitDependency(neo4j any) error {
	ctnr.KamsiCtnr = kamsi.NewContainer(nil)
	ctnr.DefineDatabase()
	ctnr.MinerRepository = repository.NewRepository(ctnr.Neo4j)
	ctnr.DefineGrpcClient()
	ctnr.HubClient = externalApi.NewHubClient("", ctnr.HubGrpcClient)
	filterMap := map[craveModel.Filter]filterBusiness.IBusiness{
		craveModel.NAME:   filterBusiness.NewNameBusiness(),
		craveModel.PERSON: filterBusiness.NewPersonBusiness(),
	}
	ctnr.FilterStrategy = filterBusiness.NewStrategy(filterMap)
	pageMap := map[craveModel.Page]pageBusiness.IBusiness{
		craveModel.NamuWiki: pageBusiness.NewNamuBusiness(pageApi.NewNamuClient()),
	}
	ctnr.PageStrategy = pageBusiness.NewStrategy(pageMap)
	ctnr.MinerService = service.NewService(ctnr.PageStrategy, ctnr.FilterStrategy, ctnr.MinerRepository, ctnr.HubClient)
	ctnr.MinerController = controller.NewController(ctnr.MinerService)
	ctnr.MinerHandler = handler.NewHandlerKamsiProxy(ctnr.KamsiCtnr.Kamsi, handler.NewHandler(ctnr.MinerController))
	return nil
}

func (ctnr *Container) DefineGrpcClient() error {
	craveGzip.RegisterGzipCompressor()

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", ctnr.Variable.Dependency.HubGrpc.Ip, ctnr.Variable.Dependency.HubGrpc.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to nginx client %w", err)
	}
	ctnr.HubGrpcClient = hubPb.NewHubClient(conn)
	return nil
}

func NewContainer(router *gin.Engine) *Container {

	ctnr := &Container{}
	ctnr.InitVariable()
	ctnr.InitDependency(nil)
	ctnr.SetRouter(router)
	return ctnr
}
