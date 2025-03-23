-- +goose Up
-- +goose StatementBegin
CREATE TABLE cart (
    user_id BIGINT NOT NULL,
    sku     INT NOT NULL CHECK (sku > 0),
    count   INT NOT NULL CHECK (count >= 0),
    PRIMARY KEY(user_id, sku)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE cart;
-- +goose StatementEnd
