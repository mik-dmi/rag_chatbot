package store

import "context"

type Storage struct {
	Queries interface {
		GetResponse(string, context.Context) (string, error)
	}
	Documents interface {
		Create([]string, context.Context) error
		Delete(string, context.Context) error
	}
	Users interface {
		CreateSession(context.Context) error
		GetChatHistory(context.Context) error
	}
}


func NewWeaviateStorage( db * ) Storage{
	return Storage{
		Documents: &DocumentsStore( db)
		Users: &UsersStore(db)
		Queries: &
		
	}
}