package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/brunoos/cnterra-controller/config"
	"github.com/brunoos/cnterra-controller/controller"
	"github.com/brunoos/cnterra-controller/db"
)

//------------------------------------------------------------------------------

func main() {
	config.Initialize()
	db.Initialize()

	r := gin.Default()

	r.GET("/files/:id", controller.GetFile)
	r.GET("/files", controller.GetAllFiles)
	r.POST("/files", controller.CreateFile)
	r.DELETE("/files/:id", controller.DeleteFile)

	r.GET("/nodes", controller.GetAllNodes)
	r.POST("/nodes", controller.CreateNode)
	r.DELETE("/nodes/:id", controller.DeleteNode)

	r.POST("/load", controller.Load)

	log.Println("[INFO] Server running...")
	r.Run()
}
