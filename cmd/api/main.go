package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/xcurvnubaim/Task-1-IS/internal/configs"
)

func main() {
	// Setup configuration
	if err := configs.Setup(); err != nil {
		panic(err)
	}

	// Setup for production
	if configs.Config.ENV_MODE == "production" {
		gin.SetMode(gin.ReleaseMode)
		fmt.Println("Production mode")
	}

	// Start the server
	r := gin.Default()

	// Setup Database
	// _, _ = database.New()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	if err := r.Run(":" + configs.Config.APP_PORT); err != nil {
		panic(err)
	}
}
