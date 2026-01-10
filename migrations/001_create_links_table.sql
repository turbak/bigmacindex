-- +goose Up
CREATE TABLE links (
                       id INTEGER PRIMARY KEY AUTOINCREMENT,
                       url TEXT NOT NULL,
                       product_name TEXT NOT NULL,
                       link_type TEXT NOT NULL,
                       price_selector TEXT NOT NULL,
                       country_code TEXT NOT NULL
);

CREATE UNIQUE INDEX idx_links_unique ON links(country_code, product_name);

INSERT INTO links (url, product_name, link_type, price_selector, country_code) VALUES
('https://vkusnotochkamenu.ru/burgers/big-hit-48567', 'Big Hit', 'html', '//strong[contains(text(),''Цена'')]/text()', 'RU');

-- +goose Down
DROP TABLE IF EXISTS links;
