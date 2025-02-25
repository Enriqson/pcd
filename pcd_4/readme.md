# Install protobuff compiler for golang

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest


protoc --go_out=./server/proto --go_opt=paths=source_relative --go-grpc_out=./server/proto --go-grpc_opt=paths=source_relative service.proto
protoc --go_out=./client/proto --go_opt=paths=source_relative --go-grpc_out=./client/proto --go-grpc_opt=paths=source_relative service.proto

