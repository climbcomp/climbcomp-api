package main

import (
	"context"
	"log"
	"net"

	"github.com/climbcomp/climbcomp-go/climbcomp"
	meta_pb "github.com/climbcomp/climbcomp-go/climbcomp/meta/v1"
	"google.golang.org/grpc"
)

func main() {
	address := "0.0.0.0:3000"
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Could not listen on %v: %v", address, err)
	}

	grpcServer := grpc.NewServer()

	meta_pb.RegisterMetaAPIServer(grpcServer, newMetaAPIServer())

	log.Println("Server starting...")
	log.Fatal(grpcServer.Serve(listener))
}

func newMetaAPIServer() *metaAPIServer {
	return &metaAPIServer{}
}

type metaAPIServer struct {
}

func (s *metaAPIServer) GetVersion(ctx context.Context, req *meta_pb.GetVersionRequest) (*meta_pb.GetVersionResponse, error) {
	resp := &meta_pb.GetVersionResponse{
		Version: climbcomp.VERSION,
	}
	return resp, nil
}
