// Simple grpc client.
// Client makes grpc get, insert, get, delete requests and prints received data from service
package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"url_shortener/internal/pkg/shortener"
)

// runs grpc client.
func main() {
	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "localhost"+":"+"50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	// run concurrent create requests
	var wg sync.WaitGroup
	const goroutineAmount = 10000
	for i := 0; i < goroutineAmount; i++ {
		wg.Add(1)
		go func(i int) {
			created, err := client.Create(ctx, &shortener.CreateUrl{Url: "https://example.com/payload" + strconv.Itoa(i)})
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(created)
			wg.Done()
		}(i)
	}
	wg.Wait()

	urlResponse, err := client.Get(ctx, &shortener.GetUrl{Url: "6LAze"})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v\n", urlResponse)

	deleted, err := client.Delete(ctx, &shortener.DeleteUrl{Url: "6LAzh"})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("removed status is " + deleted.GetStatus())

	redirect, err := client.Redirect(ctx, &shortener.RedirectUrl{Url: "6LAze"})
	if err != nil {
		fmt.Println(err)
	}
	err = exec.Command("open", redirect.GetUrl()).Start()
	if err != nil {
		fmt.Println(err)
	}
}
