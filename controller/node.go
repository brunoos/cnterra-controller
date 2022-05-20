package controller

import (
	"context"
	"net/http"

	"github.com/brunoos/cnterra-controller/db"
	"github.com/brunoos/cnterra-controller/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type formCreateNode struct {
	NodeID     int               `json:"nodeid"`
	Model      string            `json:"model"`
	Attributes map[string]string `json:"attributes"`
}

//------------------------------------------------------------------------------

func GetAllNodes(c *gin.Context) {
	col := db.DB.Collection("nodes")
	cur, err := col.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error retrieving nodes",
		})
		return
	}

	defer cur.Close(context.Background())

	nodes := make([]*model.Node, 0)
	for cur.Next(context.Background()) {
		node := new(model.Node)
		if err = cur.Decode(node); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error retrieving nodes",
			})
			return
		}
		nodes = append(nodes, node)
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

	if len(form.Model) == 0 || form.NodeID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid parameters",
		})
		return
	}

	node := model.Node{
		ID:         primitive.NewObjectID(),
		NodeID:     form.NodeID,
		Model:      form.Model,
		Attributes: form.Attributes,
	}

	col := db.DB.Collection("nodes")
	if _, err := col.InsertOne(context.Background(), &node); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error creating a new node",
		})
		return
	}

	c.JSON(http.StatusCreated, &node)
}

func DeleteNode(c *gin.Context) {
	param := c.Param("id")
	id, err := primitive.ObjectIDFromHex(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid node ID",
		})
		return
	}

	col := db.DB.Collection("nodes")
	res, err := col.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error deleting the node",
		})
		return
	}

	if res.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "node not found",
		})
		return
	}

	c.Status(http.StatusOK)
}
