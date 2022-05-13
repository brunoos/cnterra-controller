package model

type Node struct {
	BaseModel
	Nodeid     int    `json:"nodeid"`
	Model      string `json:"model"`
	Enabled    bool   `json:"enabled"`
	Attributes JSONB  `json:"attributes"`
}
