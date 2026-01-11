package main

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/turbak/bigmacindex/internal/app"
	"github.com/turbak/bigmacindex/internal/storage/links"
	"github.com/turbak/bigmacindex/internal/storage/prices"
)

func main() {
	ctx := context.Background()

	db, err := sql.Open("sqlite3", "./bigmacindex.db")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	linksRepo := links.NewRepository(db)
	pricesRepo := prices.NewRepository(db)

	linksRoutes := app.NewLinksRoutes(linksRepo)

	pricesApp := app.NewApp(linksRoutes, pricesRepo)

	if err := pricesApp.SetupRoutes(ctx); err != nil {
		log.Fatalf("failed to set up routes: %v", err)
	}
}
