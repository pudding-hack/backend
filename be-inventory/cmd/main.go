package main

import (
	"context"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/pudding-hack/backend/be-inventory/cmd/rest"
	"github.com/pudding-hack/backend/be-inventory/internal/model/history"
	"github.com/pudding-hack/backend/be-inventory/internal/model/item"
	"github.com/pudding-hack/backend/be-inventory/internal/model/unit"
	"github.com/pudding-hack/backend/be-inventory/internal/use_case"
	"github.com/pudding-hack/backend/conn"
	"github.com/pudding-hack/backend/lib"
)

func main() {
	cfg := lib.LoadConfigByFile("./cmd", "config", "yml")

	awscfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-1")) // Replace with your desired AWS region
	if err != nil {
		log.Fatalf("Failed to load AWS configuration: %v", err)
	}

	db, err := conn.NewConnectionManager(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to sql server: %v", err)
	}

	rekognitionSvc := rekognition.NewFromConfig(awscfg)

	defer db.Close()

	connQuery := db.GetQuery()

	wg := new(sync.WaitGroup)
	wg.Add(1)

	go func() {
		ctx := context.Background()
		inventoryRepository := item.New(cfg, connQuery)
		unitRepository := unit.New(cfg, connQuery)
		historyRepository := history.New(cfg, connQuery)
		inventoryService := use_case.NewService(cfg, inventoryRepository, unitRepository, historyRepository, rekognitionSvc)
		requestHandler := rest.NewHandler(inventoryService)

		err := rest.Run(ctx, *cfg, requestHandler)
		if err != nil {
			log.Fatal(err)
		}

		wg.Done()
	}()

	wg.Wait()
}
