package main

import (
	"log"

	"github.com/pudding-hack/backend/conn"
	"github.com/pudding-hack/backend/lib"
)

func main() {
	cfg := lib.LoadConfigByFile("./cmd", "config", "yml")
	log.Println(cfg)
	db, err := conn.NewConnectionManager(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to sql server: %v", err)
	}

	defer db.Close()
}
