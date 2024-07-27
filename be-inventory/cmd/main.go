package main

import (
	"context"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/pudding-hack/backend/be-inventory/cmd/rest"
	"github.com/pudding-hack/backend/be-inventory/internal/use_case"
	"github.com/pudding-hack/backend/conn"
	"github.com/pudding-hack/backend/lib"
)

func main() {
	cfg := lib.LoadConfigByFile("./cmd", "config", "yml")

	awscfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-1"))
	if err != nil {
		log.Fatalf("Failed to load AWS configuration: %v", err)
	}

	db, err := conn.NewConnectionManager(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to sql server: %v", err)
	}

	rekognitionSvc := rekognition.NewFromConfig(awscfg)

	defer db.Close()

	wg := new(sync.WaitGroup)
	wg.Add(1)

	go func() {
		ctx := context.Background()
		inventoryService := use_case.NewService(db, cfg, rekognitionSvc)
		requestHandler := rest.NewHandler(inventoryService)

		err := rest.Run(ctx, *cfg, requestHandler)
		if err != nil {
			log.Fatal(err)
		}

		wg.Done()
	}()

	wg.Wait()
}
