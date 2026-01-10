package main

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/turbak/bigmacindex/internal/poller"
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

	pricePoller := poller.NewPoller(linksRepo, pricesRepo)
	if err := pricePoller.Poll(ctx); err != nil {
		log.Fatalf("failed to poll prices: %v", err)
	}

	log.Println("Price polling completed successfully")
}
