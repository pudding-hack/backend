package main

import (
	"context"
	"log"
	"sync"

	"github.com/pudding-hack/backend/be-inventory/cmd/rest"
	"github.com/pudding-hack/backend/be-inventory/internal/model/item"
	"github.com/pudding-hack/backend/be-inventory/internal/use_case"
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

	connQuery := db.GetQuery()

	wg := new(sync.WaitGroup)
	wg.Add(1)

	go func() {
		ctx := context.Background()
		inventoryRepository := item.New(cfg, connQuery)
		inventoryService := use_case.NewService(cfg, inventoryRepository)
		requestHandler := rest.NewHandler(inventoryService)

		err := rest.Run(ctx, *cfg, requestHandler)
		if err != nil {
			log.Fatal(err)
		}

		wg.Done()
	}()

	wg.Wait()
}
