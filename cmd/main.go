package main

import (
	"xs/internal/auth"

	auth_controller "xs/internal/controller/http/v1/auth"

	user_controller "xs/internal/controller/http/v1/user"
	"xs/internal/pkg/config"
	"xs/internal/pkg/repository/postgres"

	user_repo "xs/internal/repository/postgres/user"

	"fmt"
	"log"
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

	//repository
	userRepo := user_repo.NewRepository(postgresDB)

	// mediaRepo := media_repo.NewRepository(postgresDB, fileService)

	//controller
	userController := user_controller.NewController(userRepo)
	authController := auth_controller.NewController(userRepo, authenticator)

	// router
	r := router.New(authenticator, userController, authController)
	fmt.Println(cfg.Port)
	log.Fatalln(r.Init(fmt.Sprintf(":%s", "8080")))
}
