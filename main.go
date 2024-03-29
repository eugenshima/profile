// Package main is an entry point to this microservice
package main

import (
	"context"
	"fmt"
	"net"

	cfgrtn "github.com/eugenshima/profile/internal/config"
	"github.com/eugenshima/profile/internal/handlers"
	"github.com/eugenshima/profile/internal/repository"
	"github.com/eugenshima/profile/internal/service"
	proto "github.com/eugenshima/profile/proto"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/jackc/pgx/v4/pgxpool"
)

// NewDBPsql function provides Connection with PostgreSQL database
func NewDBPsql(env string) (*pgxpool.Pool, error) {
	// Initialization a connect configuration for a PostgreSQL using pgx driver
	config, err := pgxpool.ParseConfig(env)
	if err != nil {
		return nil, fmt.Errorf("error connection to PostgreSQL: %v", err)
	}

	// Establishing a new connection to a PostgreSQL database using the pgx driver
	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("error connection to PostgreSQL: %v", err)
	}
	// Output to console
	fmt.Println("Connected to PostgreSQL!")

	return pool, nil
}

// main function of our microservice
func main() {
	cfg, err := cfgrtn.NewConfig()
	if err != nil {
		fmt.Printf("Error extracting env variables: %v", err)
		return
	}
	pool, err := NewDBPsql(cfg.PgxDBAddr)
	if err != nil {
		logrus.WithFields(logrus.Fields{"PgxDBAddr: ": cfg.PgxDBAddr}).Errorf("NewDBPsql: %v", err)
	}

	rps := repository.NewProfileRepository(pool)
	srv := service.NewProfileService(rps)
	handler := handlers.NewProfileHandler(srv)

	lis, err := net.Listen("tcp", "127.0.0.1:8082")
	if err != nil {
		logrus.Fatalf("cannot create listener: %s", err)
	}

	serverRegistrar := grpc.NewServer()
	proto.RegisterProfilesServer(serverRegistrar, handler)
	err = serverRegistrar.Serve(lis)
	if err != nil {
		logrus.Fatalf("cannot start server: %s", err)
	}
}
