

-- +goose Up
CREATE TABLE cards (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    card_desc TEXT,
    board_id INTEGER NOT NULL REFERENCES boards(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE cards;