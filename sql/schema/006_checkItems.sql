

-- +goose Up
CREATE TABLE check_items (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    check INTEGER NOT NULL,
    card_id INTEGER NOT NULL REFERENCES cards(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE check_items;