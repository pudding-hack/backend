package main

import (
	"context"
	"log"
	"sync"

	"github.com/pudding-hack/backend/be-auth/cmd/rest"
	"github.com/pudding-hack/backend/be-auth/internal/model"
	"github.com/pudding-hack/backend/be-auth/internal/use_case"
	"github.com/pudding-hack/backend/conn"
	"github.com/pudding-hack/backend/lib"
)

func main() {
	cfg := lib.LoadConfigByFile("./cmd", "config", "yml")

	db, err := conn.NewConnectionManager(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to sql server: %v", err)
	}

	defer db.Close()

	_, redisPool := lib.InitRedis(cfg.Redis)
	connQuery := db.GetQuery()

	wg := new(sync.WaitGroup)
	wg.Add(1)

	go func() {
		ctx := context.Background()
		authRepository := model.New(cfg, connQuery, redisPool)
		authService := use_case.NewService(cfg, authRepository)
		requestHandler := rest.NewHandler(authService)

		err := rest.Run(ctx, *cfg, requestHandler)
		if err != nil {
			log.Fatal(err)
		}

		wg.Done()
	}()

	wg.Wait()
}
