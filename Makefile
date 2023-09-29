proto:
	protoc --go_out=. --go-grpc_out=.  proto/auth.proto

test:
	go test -v -cover ./...