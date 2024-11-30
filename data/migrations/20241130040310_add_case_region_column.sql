-- +goose Up
-- +goose StatementBegin
ALTER TABLE cases 
    ADD COLUMN region text NOT NULL DEFAULT 'MX_DGO_DGO';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE cases
    DROP COLUMN region;
-- +goose StatementEnd
