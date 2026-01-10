-- +goose Up
CREATE TABLE prices (
                        id INTEGER PRIMARY KEY AUTOINCREMENT,
                        product_name TEXT NOT NULL,
                        price INTEGER NOT NULL,
                        price_cents INTEGER NOT NULL,
                        currency TEXT NOT NULL,
                        country_code TEXT NOT NULL,
                        created_date TEXT NOT NULL
);

CREATE UNIQUE INDEX idx_prices_unique ON prices(product_name, country_code, created_date);

-- +goose Down
DROP TABLE IF EXISTS prices;
