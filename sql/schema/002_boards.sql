

-- +goose Up
CREATE TABLE boards (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE boards;