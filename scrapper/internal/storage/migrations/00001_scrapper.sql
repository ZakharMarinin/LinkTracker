-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS users (
    chat_id INT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS links (
    link_id SERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    alias TEXT NOT NULL,
    last_update TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS user_links (
    id SERIAL PRIMARY KEY,
    chat_id INT NOT NULL REFERENCES users(chat_id),
    link_id INT NOT NULL REFERENCES links(link_id),
    alias TEXT NOT NULL,
    description TEXT
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS links CASCADE;
DROP TABLE IF EXISTS user_links CASCADE;
-- +goose StatementEnd
