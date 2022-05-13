package model

type Node struct {
	BaseModel
	NodeID     int    `json:"nodeid" gorm:"column:nodeid"`
	Model      string `json:"model"`
	Enabled    bool   `json:"enabled"`
	Attributes JSONB  `json:"attributes"`
}
