-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS cart_items (
    id SERIAL PRIMARY KEY,
    cart_id INT NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
    product_id INT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    quantity INT NOT NULL CHECK (quantity > 0),
    UNIQUE (cart_id, product_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cart_items;
-- +goose StatementEnd
