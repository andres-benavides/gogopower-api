package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

var DB *pgxpool.Pool

func InitDB(dsn string) {
	var err error
	DB, err = pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatal("Error al conectar a la base de datos:", err)
	}

	createTableSQL := `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

	CREATE TABLE IF NOT EXISTS ratings (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		ticker TEXT,
		target_from TEXT,
		target_to TEXT,
		company TEXT,
		action TEXT,
		brokerage TEXT,
		rating_from TEXT,
		rating_to TEXT,
		time TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err = DB.Exec(context.Background(), createTableSQL)
	if err != nil {
		log.Fatal("Error al crear tabla:", err)
	}
}
