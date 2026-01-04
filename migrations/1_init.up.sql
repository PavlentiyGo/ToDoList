CREATE TABLE IF NOT EXISTS tasks(
    id SERIAL PRIMARY KEY,
    title TEXT,
    description TEXT,
    done BOOLEAN DEFAULT FALSE
);