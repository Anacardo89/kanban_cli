package storage

const (
	CreateDB = `
	CREATE TABLE IF NOT EXISTS projects (
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL,
	);
	
	CREATE TABLE IF NOT EXISTS boards (
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL,
		project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE
	);
	
	CREATE TABLE IF NOT EXISTS labels (
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL,
		color TEXT NOT NULL,
		project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS cards (
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL,
		card_desc TEXT,
		board_id INTEGER NOT NULL REFERENCES boards(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS card_labels (
		card_id INTEGER NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
		label_id INTEGER NOT NULL REFERENCES labels(id) ON DELETE CASCADE,
		UNIQUE(card_id, label_id)
	);
	
	CREATE TABLE IF NOT EXISTS check_items (
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL,
		check INTEGER NOT NULL,
		card_id INTEGER NOT NULL REFERENCES cards(id) ON DELETE CASCADE
	);`
	CreateTableProjects = `
	CREATE TABLE IF NOT EXISTS projects (
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL,
	);`
	CreateTableBoards = `
	CREATE TABLE IF NOT EXISTS boards (
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL,
		project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE
	);`
	CreateTableLabels = `
	CREATE TABLE IF NOT EXISTS labels (
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL,
		color TEXT NOT NULL,
		project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE
	);`
	CreateTableCards = `
	CREATE TABLE IF NOT EXISTS cards (
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL,
		card_desc TEXT,
		board_id INTEGER NOT NULL REFERENCES boards(id) ON DELETE CASCADE
	);`
	CreateTableCardLabels = `
	CREATE TABLE IF NOT EXISTS card_labels (
		card_id INTEGER NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
		label_id INTEGER NOT NULL REFERENCES labels(id) ON DELETE CASCADE,
		UNIQUE(card_id, label_id)
	);`
	CreateTableCheckItems = `
	CREATE TABLE IF NOT EXISTS check_items (
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL,
		check INTEGER NOT NULL,
		card_id INTEGER NOT NULL REFERENCES cards(id) ON DELETE CASCADE
	);`
)
