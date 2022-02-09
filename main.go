package main

import (
	DATABASE "shopit/database"
	HELPERS "shopit/helpers"

	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

var PORT string
var GIN_MODE string

func init() {
	PORT = HELPERS.GetEnv("PORT")
	GIN_MODE = HELPERS.GetEnv("GIN_MODE")
	gin.SetMode(GIN_MODE)
}

func main() {
	router := gin.Default()
	router.Use(cors.Default())

	// Serve frontend static files
	router.Use(static.Serve("/", static.LocalFile("./public", true)))
	router.NoRoute(func(c *gin.Context) {
		c.File("./public/index.html")
	})

	api := router.Group("/api")
	{
		api.GET("/status", DATABASE.CheckConnection)

		itemApi := api.Group("/item")
		{
			itemApi.GET("/list", DATABASE.GetItemList)
			itemApi.GET("/:itemId", DATABASE.GetItemDetails)
			itemApi.POST("/", DATABASE.AddItem)
			itemApi.DELETE("/", DATABASE.DeleteItem)
			itemApi.PUT("/", DATABASE.UpdateItem)
			itemApi.POST("/image", DATABASE.UploadImage)
		}
	}

	router.Run(":" + PORT)
}
