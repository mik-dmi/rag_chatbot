package store

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

type VectorsStore struct {
	client *weaviate.Client
}

func (d *VectorsStore) CreateVectors(clientID []string, ctx context.Context) error {
	return nil
}
func (d *VectorsStore) DeleteVectors(deleteDocument string, ctx context.Context) error {
	return nil
}
