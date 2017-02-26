package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hyperboloide/static"
	"github.com/hyperboloide/static/demo/files"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	rootGroup := r.Group("/")
	{
		// Asset and AssetNames comes from the generated bindata.go file
		staticHandler := static.NewHandler(files.Asset, files.AssetNames)

		// filename to lookup for directory indexes.
		staticHandler.AddIndexes("index.html", "index.htm")

		// register the staticHandler with the group.
		staticHandler.Register(rootGroup)
	}
	r.Run(":8000")
}
