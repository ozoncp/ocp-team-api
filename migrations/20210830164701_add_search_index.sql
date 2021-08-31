-- +goose Up
-- +goose StatementBegin
ALTER TABLE team ADD COLUMN tsv tsvector;

COMMENT ON COLUMN team.tsv IS 'The column for search index depending on name and description columns';

UPDATE team SET tsv = setweight(to_tsvector(name), 'A') ||
                      setweight(to_tsvector(description), 'B') WHERE is_deleted = FALSE;

CREATE INDEX ix_team_tsv ON team USING GIN(tsv) WHERE is_deleted = FALSE;

CREATE FUNCTION trigger_tsvector_team() RETURNS TRIGGER AS $$
BEGIN
    NEW.tsv = setweight(to_tsvector(NEW.name), 'A')
                  || setweight(to_tsvector(NEW.description), 'B');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER team_tsvector_column
    BEFORE INSERT OR UPDATE ON team
    FOR EACH ROW
EXECUTE PROCEDURE trigger_tsvector_team();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX ix_team_tsv;
ALTER TABLE team DROP COLUMN tsv RESTRICT;
DROP FUNCTION trigger_tsvector_team CASCADE;
-- +goose StatementEnd
