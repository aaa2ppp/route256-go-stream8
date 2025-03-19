-- +goose Up
-- +goose StatementBegin
CREATE TABLE cart (
    user_id BIGINT NOT NULL,
    sku     INT NOT NULL,
    count   INT NOT NULL,
    PRIMARY KEY(user_id, sku)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE cart;
-- +goose StatementEnd
