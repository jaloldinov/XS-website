package main

import (
	"fmt"
	"log"
	"xs/internal/auth"
	"xs/internal/pkg/config"
	"xs/internal/pkg/repository/postgres"
	"xs/internal/router"
)

func main() {
	// config
	cfg := config.GetConf()
	fmt.Println("user:", cfg.DBUsername)
	// databases
	postgresDB := postgres.New(cfg.DBUsername, cfg.DBPassword, cfg.Port, cfg.DBName)
	// authenticator
	authenticator := auth.New(postgresDB)

	// router
	r := router.New(authenticator, postgresDB)
	fmt.Println(cfg.Port)
	log.Fatalln(r.Init(fmt.Sprintf(":%s", "8080")))

}
