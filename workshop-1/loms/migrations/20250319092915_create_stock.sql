-- +goose Up
-- +goose StatementBegin
CREATE TABLE stock (
    sku      BIGINT PRIMARY KEY,
    count    BIGINT NOT NULL,
    reserved BIGINT NOT NULL DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE stock;
-- +goose StatementEnd
