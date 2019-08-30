package cmd

import (
	"context"
	"log"

	meta_pb "github.com/climbcomp/climbcomp-go/climbcomp/meta/v1"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

func OnMetaVersionCmd(c *cli.Context) error {
	address := c.String("address")
	log.Printf("Dialing %v", address)

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not dial %v: %v", address, err)
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