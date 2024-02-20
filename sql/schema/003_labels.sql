

-- +goose Up
CREATE TABLE labels (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    color TEXT NOT NULL,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE labels;