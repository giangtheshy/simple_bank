package main

import (
	"database/sql"
	"log"

	"github.com/giangtheshy/simple_bank/api"
	db "github.com/giangtheshy/simple_bank/db/sqlc"
	"github.com/giangtheshy/simple_bank/util"
	_ "github.com/lib/pq"
)


func main() {
	config,err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.InitStore(conn)
	server,err := api.NewServer(config,store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}


	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}