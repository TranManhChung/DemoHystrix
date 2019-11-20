BIN = zpi-e-voucher
REVISION = $(shell git log | head -n 1 | cut  -f 2 -d ' ')

clean:
	rm -f $(BIN)

build: clean
	GOOS=linux go build -ldflags "-X main.revision=$(REVISION)"

genpb:
	protoc -I/usr/local/include -Igrpc-gen/admin \
		-I$$GOPATH/src \
		-I$$GOPATH/src/github.com/gogo/protobuf/protobuf \
		-I$$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--gogo_out=plugins=grpc:grpc-gen/admin \
		grpc-gen/admin/admin.proto

	protoc -I/usr/local/include -Igrpc-gen/voucher \
		-I$$GOPATH/src \
		-I$$GOPATH/src/github.com/gogo/protobuf/protobuf \
		-I$$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--gogo_out=plugins=grpc:grpc-gen/voucher \
		grpc-gen/voucher/voucher.proto

	protoc -I/usr/local/include -Igrpc-gen/cashback \
		-I$$GOPATH/src \
		-I$$GOPATH/src/github.com/gogo/protobuf/protobuf \
		-I$$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--gogo_out=plugins=grpc:grpc-gen/cashback \
		grpc-gen/cashback/cashback.proto

rsync:
	# rsync -avz zpi-voucher* root@10.30.83.2:/zserver/go-projects/zpi-voucher/
	# ssh root@10.30.83.2 sh /zserver/go-projects/zpi-voucher/runserver.sh restart
	
deploy: build rsync
	
