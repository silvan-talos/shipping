package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/silvan-talos/shipping/http"
	"github.com/silvan-talos/shipping/inmem"
	"github.com/silvan-talos/shipping/product"
)

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("failed to create listener, error:", err)
	}
	productService := product.NewService(product.ServiceArgs{
		Packs: inmem.NewPackRepository(),
	})
	server := http.NewServer(http.ServerArgs{
		ProductService: productService,
	})
	errs := make(chan error, 2)
	go func() {
		quit := make(chan os.Signal, 2)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		errs <- fmt.Errorf("signal: %s", <-quit)
	}()
	go func() {
		err = server.Serve(lis)
		log.Println("HTTP server stopped, err:", err)
		errs <- fmt.Errorf("err: %w", err)
	}()

	log.Println("exiting,", <-errs)
}
