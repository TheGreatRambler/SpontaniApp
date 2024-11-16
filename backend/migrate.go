package main

import (
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	m, err := migrate.New("file://migrations", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		panic(err)
	}
}
