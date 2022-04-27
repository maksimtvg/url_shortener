package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"os/signal"
	"syscall"
	"url_shortener/internal/app"
	"url_shortener/internal/config"
	"url_shortener/internal/database"
	"url_shortener/internal/pkg/shortener"
)

func main() {
	cfg := config.NewConfig()
	dbCfg := database.NewDBConfig()
	dbConn, err := database.Connect(dbCfg)
	if err != nil {
		log.Fatalf("%s", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("listen fatal error: %s", err)
	}

	grpcServer := grpc.NewServer()
	shortener.RegisterUrlShortenerServiceServer(grpcServer, app.NewUrlShortener(dbConn))

	go func(listener net.Listener) {
		if err = grpcServer.Serve(listener); err != nil && err != grpc.ErrServerStopped {
			log.Fatalf("server start fatal error: %s", err)
		}
	}(lis)

	<-ctx.Done()

	grpcServer.GracefulStop()
	//stopGraceful()
	log.Println("server is stopped")
}

func stopGraceful() {

}
