package util

import (
	"context"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
)

// GrpcRetryCfg ...
type GrpcRetryCfg struct {
	RetryTimeout time.Duration
	Scalar       time.Duration
	NumOfRetry   uint
}

// NewConnection ...
func NewConnection(address string, interceptors ...grpc.UnaryClientInterceptor) *grpc.ClientConn {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure(), grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(interceptors...)))
	conn, _ := grpc.Dial(address, opts...)
	return conn
}

// GrpcRetryUnaryClientInterceptor ....
func GrpcRetryUnaryClientInterceptor(grpcRetryCfg GrpcRetryCfg) grpc.UnaryClientInterceptor {
	return grpc_retry.UnaryClientInterceptor([]grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffExponential(grpcRetryCfg.Scalar)),
		grpc_retry.WithCodes(codes.Unavailable),
		grpc_retry.WithMax(grpcRetryCfg.NumOfRetry),
		grpc_retry.WithPerRetryTimeout(grpcRetryCfg.RetryTimeout),
	}...)
}

// HystrixEnableFlag is a flag for enable/disable hystrix
var HystrixEnableFlag = true

// UnaryClientInterceptor ...
func UnaryClientInterceptor(hystrixCfg hystrix.CommandConfig) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		if HystrixEnableFlag {
			hystrix.ConfigureCommand(method, hystrixCfg)
			err := hystrix.Do(method, func() (err error) {
				err = invoker(ctx, method, req, reply, cc, opts...)
				return err
			}, nil)

			return err
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
