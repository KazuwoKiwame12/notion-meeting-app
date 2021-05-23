package main

import (
	"app/client"
	"app/config"
	"app/usecase"
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("%v\n", err)
	}
}

func main() {
	processManager := map[string]context.CancelFunc{}
	parentCtx, parentCancel := context.WithCancel(context.Background())
	defer parentCancel()
	commandUC := &usecase.CommandUsecase{
		ProcessManager: processManager,
		ParentContext:  parentCtx,
	}
	s := client.NewServer(commandUC)
	go func() {
		if err := s.Start(":" + config.Port()); err != nil {
			s.Logger.Fatalf("shutting down the server with error: %v", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	s.Logger.Infof("SIGNAL %d received, then shutting dow....", <-quit)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		s.Logger.Fatal(err)
	}
}
