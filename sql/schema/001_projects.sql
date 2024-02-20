

-- +goose Up
CREATE TABLE projects (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
);

-- +goose Down
DROP TABLE projects;