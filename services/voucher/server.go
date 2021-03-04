package voucher

import (
	"context"
	"fmt"
	"time"

	admin "github.com/TranManhChung/DemoHystrix/grpc-gen/admin"
	pb "github.com/TranManhChung/DemoHystrix/grpc-gen/voucher"
	"github.com/TranManhChung/DemoHystrix/util"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	"gitlab.360live.vn/zalopay/go-common/common"
	"gitlab.360live.vn/zalopay/go-common/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Server ...
type Server struct {
	tracer      opentracing.Tracer
	logger      log.Factory
	adminClient admin.EVoucherServiceClient
}

// NewServer ...
func NewServer(tracer opentracing.Tracer, logger log.Factory) *Server {

	conn := util.NewConnection("localhost:54121")
	return &Server{
		tracer:      tracer,
		logger:      logger,
		adminClient: admin.NewEVoucherServiceClient(conn),
	}
}

// Run starts the server
func (s *Server) Run() error {
	server := common.NewGrpcServer(s.registerServer, "/evoucher/voucher", s.tracer)
	server.WithLogger(s.logger)
	server.AddShutdownHook(func() {
	})

	port := viper.GetInt("voucher.grpc_port")
	s.logger.Bg().Info("Start listen", zap.Int("port", port))
	return server.Run(port)
}

func (s *Server) registerServer(server *grpc.Server) {
	pb.RegisterEVoucherServiceServer(server, s)
}

// SayHello ...
func (s *Server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	fmt.Println("nothing")
	time.Sleep(3 * time.Second)
	_, e := s.adminClient.SayHello(ctx, &admin.HelloRequest{})
	fmt.Println(e)
	return &pb.HelloResponse{}, nil
}
