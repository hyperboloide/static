package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hyperboloide/static"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	rootGroup := r.Group("/")
	{
		staticHandler := static.NewHandler(Asset, AssetNames)
		staticHandler.AddIndexes("index.html")
		staticHandler.Register(rootGroup)
	}
	r.Run(":8000")
}
