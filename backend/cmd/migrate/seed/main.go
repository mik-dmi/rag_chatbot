package main

import (
	"log"

	"github.com/mik-dmi/rag_chatbot/backend/internal/db"
	"github.com/mik-dmi/rag_chatbot/backend/internal/env"
	"github.com/mik-dmi/rag_chatbot/backend/internal/store"
)

func main() {

	addr := env.GetString("POSTGRES_ADDR", "postgres://admin:adminpassword@localhost:5492/postgres_rag?sslmode=disable")
	conn, err := db.NewPostgreClient(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	store := store.NewPostgreStorage(conn)

	db.Seed(store)
}
