package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Node struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	NodeID     int                `json:"nodeid" bson:"nodeid"`
	Model      string             `json:"model" bson:"model"`
	Attributes map[string]string  `json:"attributes" bson:"attributes"`
}
