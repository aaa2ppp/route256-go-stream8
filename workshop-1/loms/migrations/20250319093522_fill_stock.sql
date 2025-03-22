-- +goose Up
-- +goose StatementBegin
INSERT INTO stock (sku, available) VALUES
(1076963, 10),
(1148162, 10),
(1625903, 10),
(2618151, 10),
(2956315, 10),
(2958025, 10),
(3596599, 10),
(3618852, 10),
(4288068, 10),
(4465995, 10);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM stock;
-- +goose StatementEnd
