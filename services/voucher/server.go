package voucher

import (
	"context"
	"fmt"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	"gitlab.360live.vn/zalopay/go-common/common"
	"gitlab.360live.vn/zalopay/go-common/log"
	admin "gitlab.360live.vn/zalopay/zpi-e-voucher/grpc-gen/admin"
	pb "gitlab.360live.vn/zalopay/zpi-e-voucher/grpc-gen/voucher"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
)

// Server ...
type Server struct {
	tracer      opentracing.Tracer
	logger      log.Factory
	adminClient admin.EVoucherServiceClient
}

// NewServer ...
func NewServer(tracer opentracing.Tracer, logger log.Factory) *Server {

	conn := NewConnection("localhost:54121", 3000, 3*time.Second)
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

// HystrixEnableFlag is a flag for enable/disable hystrix
var HystrixEnableFlag = true

// UnaryClientInterceptor ...
func UnaryClientInterceptor(timout int) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if HystrixEnableFlag {
			hystrix.ConfigureCommand(method,
				hystrix.CommandConfig{
					Timeout:                timout,
					MaxConcurrentRequests:  1000,
					RequestVolumeThreshold: 100,
					SleepWindow:            5000,
					ErrorPercentThreshold:  80,
				})

			err := hystrix.Do(method, func() (err error) {
				err = invoker(ctx, method, req, reply, cc, opts...)
				return err
			}, nil)

			return err
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// NewConnection ...
func NewConnection(address string, timeout int, t time.Duration) *grpc.ClientConn {
	var opts []grpc.DialOption
	optsRetry := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffExponential(50 * time.Millisecond)),
		grpc_retry.WithCodes(codes.Unavailable),
		grpc_retry.WithMax(3),
		grpc_retry.WithPerRetryTimeout(t),
	}
	opts = append(opts, grpc.WithInsecure(), grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
		UnaryClientInterceptor(timeout),
		grpc_retry.UnaryClientInterceptor(optsRetry...))))
	conn, _ := grpc.Dial(address, opts...)
	return conn
}
