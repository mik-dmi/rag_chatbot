package store

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

type Document struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
type VectorsStore struct {
	client *weaviate.Client
}

func (d *VectorsStore) CreateVectors(ctx context.Context, clientID []string) error {
	return nil
}
func (d *VectorsStore) DeleteVectors(ctx context.Context, deleteDocument string) error {
	return nil
}
