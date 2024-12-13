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

INSERT INTO cases (id, case_id, case_type, alias) VALUES ('768e3345-5e60-4f79-8c63-b2c417dcd4f9', '84/2003', 'aux2', 'Caso por defecto');
INSERT INTO cases (id, case_id, case_type, alias) VALUES ('39eb054d-bd18-4cf7-9adb-890c3f9d8ffa', '268/2013', 'civ2', 'Caso civil');
INSERT INTO cases (id, case_id, case_type, alias) VALUES ('ea70bfee-846b-4d44-b2fe-f849eb23111a', '103/2006', 'mer1', '');
INSERT INTO cases (id, case_id, case_type, alias) VALUES ('23f05122-0128-42d8-918d-72769e82c02d', '24/2024', 'fam3', null);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE cases;
-- +goose StatementEnd
