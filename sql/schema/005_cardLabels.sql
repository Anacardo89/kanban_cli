

-- +goose Up
CREATE TABLE card_labels (
	card_id INTEGER NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
	label_id INTEGER NOT NULL REFERENCES labels(id) ON DELETE CASCADE,
	UNIQUE(card_id, label_id)
);

-- +goose Down
DROP TABLE card_labels;