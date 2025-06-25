-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    hash_password TEXT NOT NULL, 
    referrer_id INTEGER NULL,
    points INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT now(),
    CONSTRAINT fk_referrer FOREIGN KEY (referrer_id) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,         
    description TEXT NOT NULL,
    reward INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE user_tasks (
    user_id INTEGER NOT NULL REFERENCES users(id),
    task_id INTEGER NOT NULL REFERENCES tasks(id),
    completed_at TIMESTAMP DEFAULT now(),
    PRIMARY KEY (user_id, task_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_tasks;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS tasks;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
