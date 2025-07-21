package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/sauravsinghs/simplebank/api"
	db "github.com/sauravsinghs/simplebank/db/sqlc"
	"github.com/sauravsinghs/simplebank/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	conn, err := sql.Open("postgres", config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	address := config.HTTPServerAddress
	if len(address) == 0 {
		address = "0.0.0.0:8080"
	}

	err = server.Start(address)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
