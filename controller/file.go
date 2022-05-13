package controller

import (
	"encoding/base64"
	"net/http"

	"github.com/brunoos/cnterra-controller/db"
	"github.com/brunoos/cnterra-controller/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type formCreateFile struct {
	Name    string
	Content string
}

//------------------------------------------------------------------------------

func GetAllFiles(c *gin.Context) {
	var files []model.File

	result := db.DB.Find(&files)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error retrieving files",
		})
		return
	}

	resp := make([]*model.FileNoContent, 0)
	for _, file := range files {
		f := new(model.FileNoContent)
		f.ID = file.ID
		f.CreatedAt = file.CreatedAt
		f.UpdatedAt = file.UpdatedAt
		f.Name = file.Name
		resp = append(resp, f)
	}

	c.JSON(http.StatusOK, gin.H{
		"files": resp,
	})
}

func GetFile(c *gin.Context) {
	var file model.File
	var err error

	id := c.Param("id")
	file.ID, err = uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid parameter",
		})
		return
	}

	result := db.DB.Find(&file)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
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

	if _, err := base64.StdEncoding.DecodeString(form.Content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid parameters",
		})
		return
	}

	file := model.File{}
	file.ID = uuid.New()
	file.Name = form.Name
	file.Content = form.Content

	result := db.DB.Create(&file)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error creating a new file",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         file.ID,
		"name":       file.Name,
		"created_at": file.CreatedAt,
		"updated_at": file.UpdatedAt,
	})
}

func DeleteFile(c *gin.Context) {
	var file model.File
	var err error

	id := c.Param("id")
	file.ID, err = uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid parameter",
		})
		return
	}

	result := db.DB.Delete(&file)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "error deleting the file",
		})
		return
	}

	c.Status(http.StatusOK)
}
