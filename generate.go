package battle

//event
//go:generate ./tools/protoc/bin/protoc --go_out=. --go-grpc_out=. event/proto/*.proto

//service
//go:generate ./tools/protoc/bin/protoc --go_out=. --go-grpc_out=. proto/*.proto
