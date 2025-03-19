-- +goose Up
-- +goose StatementBegin
CREATE TYPE order_status AS ENUM('new', 'awaiting payment', 'failed', 'payed', 'cancelled');
CREATE TABLE "order" (
    order_id BIGSERIAL PRIMARY KEY,
    user_id  BIGINT NOT NULL,
    status   order_status NOT NULL DEFAULT 'new'
);
CREATE TABLE order_items (
    order_id BIGINT NOT NULL REFERENCES "order",
    sku      BIGINT NOT NULL REFERENCES stock,
    count    INT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE order_items;
DROP TABLE "order";
DROP TYPE  order_status;
-- +goose StatementEnd
