-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "task"(
    task_id SERIAL PRIMARY KEY,
    header TEXT NOT NULL, 
    task_description TEXT NULL, 
    start_time TIMESTAMP NOT NULL, 
    end_time TIMESTAMP NOT NULL, 
    user_id INT REFERENCES "user" (user_id) ON DELETE CASCADE,
    CHECK (start_time < end_time)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "task";
-- +goose StatementEnd
