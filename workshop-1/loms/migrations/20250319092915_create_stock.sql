-- +goose Up
-- +goose StatementBegin
CREATE TABLE stock (
    sku       INT PRIMARY KEY CHECK (sku > 0),
    available BIGINT NOT NULL DEFAULT 0 CHECK (available >= 0),
    reserved  BIGINT NOT NULL DEFAULT 0 CHECK (reserved  >= 0)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE stock;
-- +goose StatementEnd
