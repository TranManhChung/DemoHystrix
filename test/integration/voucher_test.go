package integration

import (
	"context"
	"testing"

	pb "github.com/TranManhChung/DemoHystrix/grpc-gen/voucher"
	"github.com/TranManhChung/DemoHystrix/util"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/stretchr/testify/assert"
)

// go test -count=1 ./test/integration/ -run TestHytrixOnly -v
func TestHytrixOnly(t *testing.T) {
	// the processing need 3s to complete but the value of hytrix is 3s so we can not send the result to client in 3s.
	// So we get Hytrix error.
	client := pb.NewEVoucherServiceClient(util.NewConnection(
		"localhost:54111",
		util.UnaryClientInterceptor(hystrix.CommandConfig{
			Timeout:                3000, // second
			MaxConcurrentRequests:  1000,
			RequestVolumeThreshold: 100,
			SleepWindow:            5000,
			ErrorPercentThreshold:  80,
		})))
	_, e := client.SayHello(context.Background(), &pb.HelloRequest{Greeting: ""})
	assert.Equal(t, e.Error(), "hystrix: timeout")
}

// func TestSuccessCase(t *testing.T) {
// 	// hytrix == timeout == 4s . But the processing only take 3s so we do not receive the error.
// 	client := pb.NewEVoucherServiceClient(voucher.NewConnection("localhost:54111", 4000, 4*time.Second))
// 	_, e := client.SayHello(context.Background(), &pb.HelloRequest{Greeting: ""})
// 	assert.Nil(t, e)
// }

// func TestRetryCase(t *testing.T) {
// 	// retry 3 time (3*3=9s). after 9s we will get the error "DeadlineExceeded". If hytrix < 9s we will get hystrix error.
// 	client := pb.NewEVoucherServiceClient(voucher.NewConnection("localhost:54111", 10000, 3*time.Second))
// 	_, e := client.SayHello(context.Background(), &pb.HelloRequest{Greeting: ""})
// 	assert.Equal(t, e.Error(), "rpc error: code = DeadlineExceeded desc = context deadline exceeded")
// }

// grpc retry first

// hystrix first
