package admin

import (
	"context"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	"gitlab.360live.vn/zalopay/go-common/common"
	"gitlab.360live.vn/zalopay/go-common/log"
	pb "gitlab.360live.vn/zalopay/zpi-e-voucher/grpc-gen/admin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Server ...
type Server struct {
	tracer opentracing.Tracer
	logger log.Factory
}

// NewServer ...
func NewServer(tracer opentracing.Tracer, logger log.Factory) *Server {
	return &Server{
		tracer: tracer,
		logger: logger,
	}
}

// Run starts the server
func (s *Server) Run() error {
	server := common.NewGrpcServer(s.registerServer, "/evoucher/admin", s.tracer)
	server.WithLogger(s.logger)
	server.AddShutdownHook(func() {
	})

	port := viper.GetInt("admin.grpc_port")
	s.logger.Bg().Info("Start listen", zap.Int("port", port))
	return server.Run(port)
}

func (s *Server) registerServer(server *grpc.Server) {
	pb.RegisterEVoucherServiceServer(server, s)
}

// SayHello ...
func (s *Server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{}, nil
}
