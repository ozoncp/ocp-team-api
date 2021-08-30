-- +goose Up
-- +goose StatementBegin
ALTER TABLE team ADD COLUMN is_deleted BOOLEAN NOT NULL DEFAULT FALSE;

COMMENT ON COLUMN team.is_deleted IS 'The deleted flag of a team';

CREATE FUNCTION trigger_soft_delete_team() RETURNS TRIGGER AS $$
    BEGIN
        UPDATE team SET is_deleted = TRUE WHERE id = OLD.id;
        RETURN NULL;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER team_soft_delete
    BEFORE DELETE ON team
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_soft_delete_team();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE team DROP COLUMN is_deleted RESTRICT;
DROP FUNCTION trigger_soft_delete_team CASCADE;
-- +goose StatementEnd
