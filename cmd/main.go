package main

import (
	"log"

	"github.com/omarelweshy/EcomMaster-user-service/internal/logger"
	"github.com/omarelweshy/EcomMaster-user-service/internal/router"
)

func main() {

	logger.InitLogger()

	r := router.SetupRouter()
	if err := r.Run(":8000"); err != nil {
		log.Fatal(err)
	}

}
