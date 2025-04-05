package main

import (
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "streaming-service/src/routes"
    "streaming-service/src/services"
)

func main() {
    // Load environment variables
    if err := godotenv.Load(); err != nil {
        panic("Error loading .env file")
    }

    // Initialize database
    InitDB()

    // Initialize Kafka service
    kafkaService, err := services.NewKafkaService()
    if err != nil {
        panic(err)
    }

    // Initialize routes with dependencies
    routes.InitUserRoutes(DB)
    routes.InitVideoRoutes(DB, kafkaService)

    router := gin.Default()
    routes.SetupUserRoutes(router)
    routes.SetupVideoRoutes(router)

    // Start server
    router.Run(":8080")
}