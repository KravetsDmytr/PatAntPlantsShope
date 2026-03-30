-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    price NUMERIC(12,2) NOT NULL CHECK (price >= 0),
    image_url TEXT NOT NULL,
    category_id INT NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS products;
-- +goose StatementEnd
