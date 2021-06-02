package main

import (
	"app/client"
	"app/config"
	"app/domain/function"
	"app/infrastructure"
	"app/observation"
	"app/usecase"
	"context"
	"fmt"
	"log"
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
	// infrastructure初期化
	sh, err := infrastructure.NewSqlHandler()
	if err != nil {
		log.Fatalf("newSqlHandler err: %+v", err)
	}
	defer func() {
		if err := sh.DB.Close(); err != nil {
			log.Fatalf("closed err: %+v", err)
		}
	}()

	// domain初期化
	dbOp := &function.DatabaseOperater{
		SqlHandler: sh,
	}

	// usecase初期化
	processManager := make(map[string]chan<- struct{})
	defer func() {
		for _, process := raneg processManager {
			close(process)
		}
	}()
	commandUC := &usecase.CommandUsecase{
		ProcessManager: processManager,
		DBOperator: dbOp,
	}
	authorizationUC := &usecase.AuthorizationUsecase{
		DBOperator: dbOp,
	}
	slackUC := &usecase.SlackUsecase{
		DBOperator: dbOp,
	}

	// client初期化
	s := client.NewServer(commandUC, authorizationUC, slackUC)
	go func() {
		// サーバ起動
		if err := s.Start(":" + config.Port()); err != nil {
			s.Logger.Fatalf("shutting down the server with error: %v", err)
		}
	}()

	// slackの最新情報を維持し続ける
	o := &observation.SlackObserver{
		DBOperator: dbOp,
	}
	go o.KeepSatate()

	// 終了処理
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	s.Logger.Infof("SIGNAL %d received, then shutting dow....", <-quit)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		s.Logger.Fatal(err)
	}
}
