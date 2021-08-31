-- +goose Up
-- +goose StatementBegin
CREATE TABLE team(
    id  SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL
);

COMMENT ON COLUMN team.id IS 'The ID of team';
COMMENT ON COLUMN team.name IS 'The name of team';
COMMENT ON COLUMN team.description IS 'The description of the team';
-- +goose StatementEnd



-- +goose Down
-- +goose StatementBegin
DROP TABLE team;
-- +goose StatementEnd
