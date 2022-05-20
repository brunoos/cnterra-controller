package controller

import (
	"context"
	"encoding/base64"
	"net/http"

	"github.com/brunoos/cnterra-controller/db"
	"github.com/brunoos/cnterra-controller/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type formCreateFile struct {
	Name    string
	Content string
}

//------------------------------------------------------------------------------

func GetAllFiles(c *gin.Context) {
	col := db.DB.Collection("files")
	cur, err := col.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error retrieving files",
		})
		return
	}

	defer cur.Close(context.Background())

	files := make([]*model.FileNoContent, 0)
	for cur.Next(context.Background()) {
		file := new(model.FileNoContent)
		if err = cur.Decode(file); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error retrieving nodes",
			})
			return
		}
		files = append(files, file)
	}

	c.JSON(http.StatusOK, gin.H{
		"files": files,
	})
}

func GetFile(c *gin.Context) {
	param := c.Param("id")
	id, err := primitive.ObjectIDFromHex(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid file ID",
		})
		return
	}

	col := db.DB.Collection("files")
	res := col.FindOne(context.Background(), bson.M{"_id": id})
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

	var file model.File
	if err = res.Decode(&file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error retrieving the file",
		})
		return
	}

	c.JSON(http.StatusOK, &file)
}

func CreateFile(c *gin.Context) {
	form := formCreateFile{}

	if err := c.Bind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid parameters",
		})
		return
	}

	if len(form.Name) == 0 || len(form.Content) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid parameters",
		})
		return
	}

	if _, err := base64.StdEncoding.DecodeString(form.Content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid parameters",
		})
		return
	}

	file := model.File{
		ID:      primitive.NewObjectID(),
		Name:    form.Name,
		Content: form.Content,
	}

	col := db.DB.Collection("files")
	if _, err := col.InsertOne(context.Background(), &file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error creating a new file",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":   file.ID,
		"name": file.Name,
	})
}

func DeleteFile(c *gin.Context) {
	param := c.Param("id")
	id, err := primitive.ObjectIDFromHex(param)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid file ID",
		})
		return
	}

	col := db.DB.Collection("files")
	res, err := col.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error deleting the file",
		})
		return
	}

	if res.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "file not found",
		})
		return
	}

	c.Status(http.StatusOK)
}
