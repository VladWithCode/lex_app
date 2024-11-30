-- +goose Up
-- +goose StatementBegin
CREATE TABLE accords (
    id text PRIMARY KEY NOT NULL,
    for_case text NOT NULL,
    content text NOT NULL,
    date_recovered real NOT NULL,

    FOREIGN KEY (for_case) REFERENCES cases(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE accords;
-- +goose StatementEnd
