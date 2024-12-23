-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX accords_unique_case_date_idx ON accords (for_case, date);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX accords_unique_case_date_idx;
-- +goose StatementEnd
