package cmd

import (
	"context"

	"github.com/climbcomp/climbcomp-api/conf"
	meta_pb "github.com/climbcomp/climbcomp-go/climbcomp/meta/v1"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

func OnMetaVersionCmd(c *cli.Context) error {
	config := conf.Instance()
	log.Printf("Dialing %v", config.Address)

	conn, err := grpc.Dial(config.Address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not dial %v: %v", config.Address, err)
	}
	defer conn.Close()

	client := meta_pb.NewMetaAPIClient(conn)

	ctx := context.Background()
	request := &meta_pb.GetVersionRequest{}
	response, err := client.GetVersion(ctx, request)
	if err != nil {
		log.Fatalf("Error: %v - %v", err, response)
	}

	log.Printf("Server version: %v", response.Version)

	return nil
}
