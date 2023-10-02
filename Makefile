proto:
	protoc --go_out=. --go-grpc_out=. --grpc-gateway_out=.  ./api/proto/auth.proto 

test:
	go test -v -cover ./...

db:
	docker start postgres15

run:
	go run ./cmd/server/main.go
