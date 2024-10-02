CREATE TABLE IF NOT EXISTS songs (
    id SERIAL PRIMARY KEY,
    "group" VARCHAR(255) NOT NULL, -- Используем двойные кавычки для зарезервированного слова
    title VARCHAR(255) NOT NULL,
    release_date DATE,
    text TEXT,
    link VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
