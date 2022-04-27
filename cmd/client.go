package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
	"url_shortener/internal/config"
	"url_shortener/internal/pkg/shortener"
)

const RequestTimeout = 5

func main() {
	cfg := config.NewConfig()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*RequestTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, cfg.Host+":"+cfg.GRPCPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("client start error: %s", err)
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("client close error: %s", err)
		}
	}(conn)

	client := shortener.NewUrlShortenerServiceClient(conn)
	_ = client
	//TODO
}
