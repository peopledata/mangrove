package schema

type NftData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	DidDoc      string `json:"image"` // 该字段用来存放did的doc文档
}
