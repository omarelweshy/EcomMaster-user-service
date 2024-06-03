package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/omarelweshy/EcomMaster-user-service/internal/handler"
	"github.com/omarelweshy/EcomMaster-user-service/internal/model"
	"github.com/omarelweshy/EcomMaster-user-service/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=omarelweshy dbname=ecommaster_user_service port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&model.User{})
	userService := &service.UserService{DB: db}
	userHandler := &handler.UserHandler{UserService: userService}

	r := gin.Default()
	r.POST("/register", userHandler.Register)

	if err := r.Run(":8000"); err != nil {
		log.Fatal(err)
	}

}
