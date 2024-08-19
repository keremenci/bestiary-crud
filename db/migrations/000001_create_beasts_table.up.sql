CREATE TABLE IF NOT EXISTS beasts (
    beast_name TEXT PRIMARY KEY,
    type TEXT NOT NULL,
    cr TEXT NOT NULL,
    attributes JSONB,
    description TEXT
);