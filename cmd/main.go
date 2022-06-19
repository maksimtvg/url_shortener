// GRPC server
package main

import (
	"url_shortener/internal/server"
)

// Runs Server.
func main() {
	s := server.NewServer()
	s.Start()
}
