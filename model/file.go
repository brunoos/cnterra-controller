package model

type File struct {
	BaseModel
	Name    string `json:"name"`
	Content string `json:"content"`
}

type FileNoContent struct {
	BaseModel
	Name string `json:"name"`
}
