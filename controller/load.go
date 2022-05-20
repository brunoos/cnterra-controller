package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/brunoos/cnterra-controller/db"
	"github.com/brunoos/cnterra-controller/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LoadConfig struct {
	Node *model.Node
	File *model.File
}

type FormLoadConfig struct {
	Node string `json:"node"`
	File string `json:"file"`
}

type FormLoad struct {
	Config []FormLoadConfig `json:"config"`
}

type Request struct {
	Content string `json:"content"`
}

//------------------------------------------------------------------------------

func loader(r chan bool, node *model.Node, file *model.File) {
	url := "http://" + node.Attributes["loader-address"] + ":" + node.Attributes["loader-port"] + "/load"

	req := Request{Content: file.Content}
	data, err := json.Marshal(&req)
	if err != nil {
		log.Printf("[ERRO] Error encoding JSON request: %s", err)
		r <- false
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Printf("[ERRO] Loader request error: %s", err)
		r <- false
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("[ERRO] Loader failure")
		r <- false
		return
	}

	r <- true
}

func Load(c *gin.Context) {
	form := FormLoad{}
	if err := c.Bind(&form); err != nil {
		log.Println("[ERRO] Error parsing JSON body")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid parameter",
		})
		return
	}

	if len(form.Config) == 0 {
		log.Println("[ERRO] Empty configuration")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid parameter",
		})
		return
	}

	config := make([]*LoadConfig, 0)
	for _, conf := range form.Config {
		node := new(model.Node)

		id, err := primitive.ObjectIDFromHex(conf.Node)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid node ID",
			})
			return
		}

		col := db.DB.Collection("nodes")
		res := col.FindOne(context.Background(), bson.M{"_id": id})
		switch res.Err() {
		case nil:
			// do nothing
		case mongo.ErrNoDocuments:
			c.JSON(http.StatusNotFound, gin.H{
				"message": "node not found",
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error retrieving the node",
			})
			return
		}

		if err = res.Decode(node); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error retrieving the node",
			})
			return
		}

		file := new(model.File)

		id, err = primitive.ObjectIDFromHex(conf.File)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid file ID",
			})
			return
		}

		col = db.DB.Collection("files")
		res = col.FindOne(context.Background(), bson.M{"_id": id})
		switch res.Err() {
		case nil:
			// do nothing
		case mongo.ErrNoDocuments:
			c.JSON(http.StatusNotFound, gin.H{
				"message": "file not found",
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error retrieving the file",
			})
			return
		}

		if err = res.Decode(file); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error retrieving the file",
			})
			return
		}

		config = append(config, &LoadConfig{
			Node: node,
			File: file,
		})
	}

	resp := make(chan bool)
	for _, conf := range config {
		go loader(resp, conf.Node, conf.File)
	}

	succ := true
	for i := 0; i < len(config); i++ {
		succ = (<-resp) && succ
	}

	if succ {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusInternalServerError)
	}
}
