package main

import (
	"shopit/database"
	"shopit/helpers"
	"shopit/middleware"

	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

var PORT string
var GIN_MODE string
var ENV string

func init() {
	PORT = helpers.GetEnv("PORT")
	ENV = helpers.GetEnv("ENV")
	GIN_MODE = helpers.GetEnv("GIN_MODE")
	gin.SetMode(GIN_MODE)
}

func main() {
	router := gin.Default()

	if ENV != "dev" {
		// Serve frontend static files
		router.Use(static.Serve("/", static.LocalFile("./public", true)))
		router.NoRoute(func(c *gin.Context) {
			c.File("./public/index.html")
		})
		router.Use(cors.New(cors.Config{AllowAllOrigins: false}))
	} else {
		router.Use(middleware.Cors(middleware.CorsConfig{IsDevMode: true}))
	}

	api := router.Group("/api")
	{
		itemApi := api.Group("/item")
		{
			itemApi.GET("/list", database.GetItemList)
			itemApi.GET("/:itemId", database.GetItemDetails)
			itemApi.POST("/", database.AddItem)
			itemApi.DELETE("/", database.DeleteItem)
			itemApi.PUT("/", database.UpdateItem)
			itemApi.POST("/image", database.UploadImage)
		}
	}

	router.Run(":" + PORT)
}
