package store

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

type Documents struct {
	UserID    string     `json:"userid"`
	Document  []Document `json:"document"`
	CreatedAt string     `json:"created_at"`
	UpdatedAt string     `json:"updated_at"`
}
type Document struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
type VectorsStore struct {
	client *weaviate.Client
}

func (d *VectorsStore) CreateVectors(ctx context.Context, documents *Documents) error {
	return nil
}
func (d *VectorsStore) DeleteVectors(ctx context.Context, deleteDocument string) error {
	return nil
}
