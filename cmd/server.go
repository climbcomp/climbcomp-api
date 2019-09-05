package cmd

import (
	"database/sql"
	"net"

	"github.com/climbcomp/climbcomp-api/conf"
	"github.com/climbcomp/climbcomp-api/meta"
	meta_pb "github.com/climbcomp/climbcomp-go/climbcomp/meta/v1"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

func OnServerCmd(c *cli.Context) error {
	config := conf.Instance()

	log.Info("Connecting to db")
	db, err := sql.Open("postgres", config.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.New(config.MigrationsUrl, config.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}

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
