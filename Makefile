.PHONY: prep protobuf server client build run_daemon run_server clean

all: build

protobuf_deps:
	@wget "https://github.com/protocolbuffers/protobuf/releases/download/v3.13.0/protoc-3.13.0-linux-x86_64.zip" -qO /tmp/protoc.zip
	@mkdir -p /tmp/protoc
	@unzip /tmp/protoc.zip -d /tmp/protoc
	@rm -f /tmp/protoc.zip
	@sudo rm -rf /usr/local/include/google/protobuf
	@sudo mkdir -p /usr/local/include/google/
	@sudo mv /tmp/protoc/include/google/protobuf /usr/local/include/google/protobuf
	@sudo mv /tmp/protoc/bin/protoc /usr/local/bin/protoc
	@rm -rf /tmp/protoc

k3c_deps:
	@sudo mkdir -p /var/lib/cni/results
	@sudo wget https://github.com/rancher/k3c/releases/download/v0.2.1/k3c-linux-amd64 -Oq /usr/local/bin/k3c
	@sudo chmod +x /usr/local/bin/k3c

grpc_deps:
	@go get -u google.golang.org/grpc
	@go get -u github.com/golang/protobuf/protoc-gen-go

prep: k3c_deps grpc_deps protobuf_deps
	@echo "Setting up your local machine."

protobuf:
	@protoc -I pkg/proto/containerd/ pkg/proto/containerd/containerd.proto --go_out=plugins=grpc:pkg/proto/containerd/
 
server:
	@go build -race -ldflags "-s -w" -o bin/server server/main.go
 
client:
	@go build -race -ldflags "-s -w" -o bin/client client/main.go

build: protobuf server client
	@echo "Building project..."

run_daemon:
	@sudo /usr/local/bin/k3c daemon

run_server: protobuf server
	@sudo bin/server

clean:
	@rm -f bin/client
	@rm -f bin/server
	@rm -f pkg/proto/containerd/containerd.pb.go
