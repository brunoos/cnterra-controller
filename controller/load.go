package controller

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/brunoos/cnterra-controller/db"
	"github.com/brunoos/cnterra-controller/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LoadConfig struct {
	Node *model.Node
	File *model.File
}

type FormLoadConfig struct {
	Node uuid.UUID `json:"node"`
	File uuid.UUID `json:"file"`
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
		result := db.DB.Where("id = ? AND enabled = ?", conf.Node, true).Find(node)
		switch result.Error {
		case nil:
			// do nothing
		case gorm.ErrRecordNotFound:
			log.Printf("[ERRO] Node '%s' not found", conf.Node)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid parameter",
			})
			return
		default:
			log.Printf("[ERRO] Error retrieving node '%s'", conf.Node)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error retrieving node information",
			})
			return
		}

		file := new(model.File)
		result = db.DB.Where("id = ?", conf.File.String()).Find(file)
		switch result.Error {
		case nil:
			// do nothing
		case gorm.ErrRecordNotFound:
			log.Printf("[ERRO] File '%s' not found", conf.File)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid parameter",
			})
			return
		default:
			log.Printf("[ERRO] Error retrieving file '%s'", conf.File)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error retrieving file information",
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
