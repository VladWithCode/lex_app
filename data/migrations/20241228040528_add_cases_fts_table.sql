-- +goose Up
-- +goose StatementBegin
 CREATE VIRTUAL TABLE cases_fts USING fts5(uuid UNINDEXED, case_id, case_type, alias, nature);

INSERT INTO cases_fts (uuid, case_id, case_type, alias, nature)
	SELECT
		id,
		case_id,
		case_type,
		alias,
		nature
    FROM cases;

CREATE TRIGGER fts_after_update_cases AFTER UPDATE ON cases
BEGIN
    INSERT INTO cases_fts (uuid, case_id, case_type, alias, nature)
        VALUES (new.id, new.case_id, new.case_type, new.alias, new.nature);
    DELETE FROM cases_fts WHERE uuid = old.id;
END;

CREATE TRIGGER fts_after_delete_cases AFTER DELETE ON cases
BEGIN
    DELETE FROM cases_fts WHERE uuid = old.id;
END;

CREATE TRIGGER fts_after_insert_cases AFTER INSERT ON cases
BEGIN
    INSERT INTO cases_fts (uuid, case_id, case_type, alias, nature)
        VALUES (new.id, new.case_id, new.case_type, new.alias, new.nature);
END;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE cases_fts;

DROP TRIGGER fts_after_update_cases ON cases;
DROP TRIGGER fts_after_insert_cases ON cases;
DROP TRIGGER fts_after_delete_cases ON cases;
-- +goose StatementEnd
