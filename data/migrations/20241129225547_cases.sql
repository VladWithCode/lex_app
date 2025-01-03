-- +goose Up
-- +goose StatementBegin
CREATE TABLE cases (
    id TEXT PRIMARY KEY NOT NULL,
    case_id TEXT NOT NULL,
    case_type TEXT NOT NULL,
    case_year TEXT DEFAULT 0,
    case_no TEXT DEFAULT 0,
    alias TEXT DEFAULT '',
    other_ids TEXT DEFAULT '',

    CONSTRAINT unique_case_id UNIQUE(case_id, case_type)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE cases;
-- +goose StatementEnd
