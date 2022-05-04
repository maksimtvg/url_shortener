// Package server.
package server

import (
	"context"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"log"
	"net"
	"os/signal"
	"syscall"
	"url_shortener/internal/config"
	"url_shortener/internal/database"
	"url_shortener/internal/pkg/shortener"
	"url_shortener/internal/repositories"
	"url_shortener/internal/services"
)

const migrationLocation = "file:///db/migrations/"

type Server struct {
	appConfig *config.Config
	dbConfig  *database.DBConfig
}

// NewServer constructs app Server
func NewServer() *Server {
	appC := config.NewConfig()
	dbC := database.NewDBConfig()

	return &Server{
		appConfig: appC,
		dbConfig:  dbC,
	}
}

// Start launches app server.
// Server starts in goroutine with context.
// In order to stop server gracefully  there are two syscall - syscall.SIGTERM, syscall.SIGINT
func (s *Server) Start() {
	s.migrate()

	dbConn, err := database.Connect(s.dbConfig)
	if err != nil {
		log.Fatalf("%s", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	lis, err := net.Listen("tcp", ":"+s.appConfig.GRPCPort)
	if err != nil {
		log.Fatalf("listen fatal error: %s", err)
	}

	grpcServer := grpc.NewServer()
	shortener.RegisterUrlShortenerServiceServer(
		grpcServer,
		services.NewUrlShortener(repositories.NewDBRepository(dbConn)),
	)
	log.Println("Running")

	go func(listener net.Listener) {
		if err = grpcServer.Serve(listener); err != nil && err != grpc.ErrServerStopped {
			log.Fatalf("server start fatal error: %s", err)
		}
	}(lis)

	<-ctx.Done()
	s.stopGracefully(grpcServer, dbConn)
}

func (s *Server) migrate() {
	m, err := migrate.New(
		migrationLocation,
		database.PgString(s.dbConfig),
	)
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
}

// stopGraceful stops server and DB gracefully
func (s *Server) stopGracefully(server *grpc.Server, dbConn *pgxpool.Pool) {
	dbConn.Close()
	server.GracefulStop()
	log.Println("server is stopped")
}
