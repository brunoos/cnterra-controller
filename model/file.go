package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type File struct {
	ID      primitive.ObjectID `json:"id" bson:"_id"`
	Name    string             `json:"name" bson:"name"`
	Content string             `json:"content" bson:"content"`
}

type FileNoContent struct {
	ID   primitive.ObjectID `json:"id" bson:"_id"`
	Name string             `json:"name" bson:"name"`
}
