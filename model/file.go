package model

import "github.com/google/uuid"

type File struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Content string    `json:"content"`
}

type FileNoContent struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
