-- +goose Up
-- +goose StatementBegin
ALTER TABLE team ADD COLUMN is_deleted BOOLEAN NOT NULL DEFAULT FALSE;

COMMENT ON COLUMN team.is_deleted IS 'The deleted flag of a team';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE team DROP COLUMN is_deleted RESTRICT;
-- +goose StatementEnd
