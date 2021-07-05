package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arfan21/getprint-order/configs"
)

func Start() error {
	mysqlClient, err := configs.NewClient()
	if err != nil {
		return err
	}

	defer mysqlClient.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	httpServer := NewRouter(mysqlClient)

	// Start server
	idleConnsClosed := make(chan struct{})
	go func() {
		defer close(idleConnsClosed)

		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)

		<-sigint
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(ctx); err != nil {
			httpServer.Logger.Fatal(err)
		}
	}()

	if err := httpServer.Start(fmt.Sprintf(":%s", port)); err != nil && err != http.ErrServerClosed {
		return err
	}
	<-idleConnsClosed
	httpServer.Logger.Print("stopped server gracefully")
	return nil
}
