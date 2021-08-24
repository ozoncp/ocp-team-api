-- +goose Up
-- +goose StatementBegin
CREATE TABLE team(
    id  SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(255) NOT NULL
)
-- +goose StatementEnd



-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS team
-- +goose StatementEnd
