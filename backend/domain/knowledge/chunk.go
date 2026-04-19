package knowledge

type Chunk struct {
	ID         string `json:"id"`
	DocumentID string `json:"document_id"`
	Content    string `json:"content"`
	Index      int    `json:"index"`
}

func NewChunk(id, documentID, content string, index int) *Chunk {
	return &Chunk{
		ID:         id,
		DocumentID: documentID,
		Content:    content,
		Index:      index,
	}
}
