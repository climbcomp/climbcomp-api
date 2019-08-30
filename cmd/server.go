package cmd

import (
	"log"
	"net"

	"github.com/climbcomp/climbcomp-api/meta"

	meta_pb "github.com/climbcomp/climbcomp-go/climbcomp/meta/v1"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

func OnServerCmd(c *cli.Context) error {
	address := "0.0.0.0:3000"
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Could not listen on %v: %v", address, err)
	}

	grpcServer := grpc.NewServer()
	meta_pb.RegisterMetaAPIServer(grpcServer, meta.NewMetaServer())

	log.Println("Server starting...")
	log.Fatal(grpcServer.Serve(listener))

	return nil
}
