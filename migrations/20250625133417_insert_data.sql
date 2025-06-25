-- +goose Up
-- +goose StatementBegin
INSERT INTO tasks (name, description, reward) VALUES
('Регистрация','Пройдите регистрацию на сайте', 50),
('Тестовое задание DeNet', 'Выполни тестовое задание в компанию DeNet', 100000);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
