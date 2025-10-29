package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jchiu21/streamsense/server/controllers"
	"github.com/jchiu21/streamsense/server/database"
)

func main() {
    if err := database.InitializeDB(); err != nil {
        log.Fatal("Failed to initialize database: ", err)
    }

    router := gin.Default()

    router.GET("/movies", controllers.GetMovies())

    if err := router.Run(":8080"); err != nil {
        fmt.Println("Failed to start server", err)
    }
} 