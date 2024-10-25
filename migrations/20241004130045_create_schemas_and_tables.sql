-- +goose Up
-- SQL запросы для создания схем и таблиц

-- Таблица authors в схеме library
CREATE TABLE IF NOT EXISTS authors
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

-- Таблица books в схеме library
CREATE TABLE IF NOT EXISTS books
(
    id        SERIAL PRIMARY KEY,
    title     VARCHAR(200) NOT NULL,
    author_id INT REFERENCES authors (id) ON DELETE CASCADE
);


-- +goose Down
-- SQL запросы для отката миграции

DROP TABLE IF EXISTS books;
DROP TABLE IF EXISTS authors;



