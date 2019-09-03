package cmd

import (
	"net"

	"github.com/climbcomp/climbcomp-api/meta"
	meta_pb "github.com/climbcomp/climbcomp-go/climbcomp/meta/v1"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

func OnServerCmd(c *cli.Context) error {
	address := "0.0.0.0:3000"
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Could not listen on %v: %v", address, err)
	}

	logger := log.StandardLogger()
	entry := log.NewEntry(logger)

	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_logrus.UnaryServerInterceptor(entry),
		),
	)
	meta_pb.RegisterMetaAPIServer(grpcServer, meta.NewMetaServer())

	log.Println("Server starting...")
	log.Fatal(grpcServer.Serve(listener))

	return nil
}
