package router

import (
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/omarelweshy/EcomMaster-user-service/docs"
	"github.com/omarelweshy/EcomMaster-user-service/internal/handler"
	"github.com/omarelweshy/EcomMaster-user-service/internal/middleware"
	"github.com/omarelweshy/EcomMaster-user-service/internal/model"
	"github.com/omarelweshy/EcomMaster-user-service/internal/repository"
	"github.com/omarelweshy/EcomMaster-user-service/internal/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupRouter() *gin.Engine {
	dsn := "host=localhost user=postgres password=omarelweshy dbname=ecommaster_user_service port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&model.User{})
	userRepository := &repository.UserRepository{DB: db}
	userService := &service.UserService{Repo: *userRepository}
	userHandler := &handler.UserHandler{UserService: userService}

	r := gin.Default()

	r.Use(middleware.Logging())
	r.Use(middleware.ErrorHandler())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/register", userHandler.Register)
	r.POST("/login", userHandler.Login)

	auth := r.Group("/auth")
	auth.Use(middleware.JWTAuthMiddleware())
	{
		auth.GET("/profile", userHandler.Profile)
		auth.PUT("/profile", userHandler.UpdateProfile)
	}

	return r
}
