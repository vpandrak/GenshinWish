package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"rest-todo/internal/Repository"
	"rest-todo/internal/server"
	"rest-todo/internal/store"
)

var (
	version = "0.1 alpha"
	build   = 10
)

func main() {
	fmt.Println("Genshin wish")
	fmt.Printf("Version %s build %v\n", version, build)

	db, err := sql.Open("pgx", "postgres://postgres:vpandrak@localhost:5432/postgres")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	repository := Repository.NewPostgresRepository(db)
	ctx := context.Background()
	store.RunRepositoryDemo(ctx, repository)
	fmt.Println("Server running on port :80")
	server.Serve(repository, ctx)

}
