-- +goose Up
-- +goose StatementBegin
ALTER TABLE cases
    ADD COLUMN nature text DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE cases
    DROP COLUMN nature;
-- +goose StatementEnd
