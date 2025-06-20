package main

import (
	"log"

	"github.com/brythnl/coop/api"
	"github.com/brythnl/coop/db"
	"github.com/brythnl/coop/db/sqlc"
	"github.com/brythnl/coop/util"
)

func main() {
	config, err := util.LoadConfig()
	if err != nil {
		log.Fatal("failed loading config:", err)
	}

	pool, err := db.ConnectDB(config.DBUrl)
	if err != nil {
		log.Fatal("failed connecting to db:", err)
	}
	defer pool.Close()

	store := sqlc.NewStore(pool)

	router := api.SetupRouter(store, config.JWTSecretKey)

	log.Printf("server starting on %s", config.ServerAddr)
	if err := router.Run(config.ServerAddr); err != nil {
		log.Fatal("server failed to start:", err)
	}
	log.Printf("server running on %s", config.ServerAddr)
}
