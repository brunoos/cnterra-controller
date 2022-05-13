package controller

import (
	"net/http"

	"github.com/brunoos/cnterra-controller/db"
	"github.com/brunoos/cnterra-controller/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type formCreateNode struct {
	Nodeid     int         `json:"nodeid"`
	Model      string      `json:"model"`
	Enabled    bool        `json:"enabled"`
	Attributes model.JSONB `json:"attributes"`
}

//------------------------------------------------------------------------------

func GetAllNodes(c *gin.Context) {
	var nodes []model.Node

	result := db.DB.Find(&nodes)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error retrieving nodes",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"nodes": nodes,
	})
}

func CreateNode(c *gin.Context) {
	form := formCreateNode{}

	if err := c.Bind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid parameters",
		})
		return
	}

	node := model.Node{}
	node.ID = uuid.New()
	node.Nodeid = form.Nodeid
	node.Model = form.Model
	node.Enabled = form.Enabled
	node.Attributes = form.Attributes

	result := db.DB.Create(&node)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error creating a new node",
		})
		return
	}

	c.JSON(http.StatusCreated, &node)
}

func DeleteNode(c *gin.Context) {
	var node model.Node
	var err error

	id := c.Param("id")
	node.ID, err = uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid parameter",
		})
		return
	}

	result := db.DB.Delete(&node)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error deleting the node",
		})
		return
	}

	c.Status(http.StatusOK)
}
