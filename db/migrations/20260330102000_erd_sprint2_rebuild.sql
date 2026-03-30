-- +goose Up
-- +goose StatementBegin
-- Пересобираем таблицы под ERD (Sprint 2), сохраняя endpoint-контракт:
-- products => name/description/image_url/price,
-- cart_products => количество + расчёт общей суммы.

-- USERS
DROP TABLE IF EXISTS cart_products;
DROP TABLE IF EXISTS cart_items;

DROP TABLE IF EXISTS carts;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS users;

-- Создаем USERS как на диаграмме (+ минимальные поля для auth).
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(255) UNIQUE NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(255) DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_login TIMESTAMP,
    password_hash VARCHAR(512) NOT NULL,
    salt VARCHAR(64) NOT NULL DEFAULT '',
    role VARCHAR(64) NOT NULL DEFAULT 'user'
);

-- CARTS
-- CATEGORIES (оставляем как нужно для Sprint 2 UI)
-- (На ERD может быть не показана, но фронт использует /categories)
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(150) UNIQUE NOT NULL
);

-- PRODUCTS как на диаграмме + дополнительные поля, нужные фронту.
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    cost NUMERIC(10,2) NOT NULL CHECK (cost >= 0),
    producer VARCHAR(255) NOT NULL DEFAULT '',
    type VARCHAR(255) NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    guarantee TIMESTAMP,
    weight NUMERIC(10,2),
    unit VARCHAR(10),
    title TEXT,

    -- для текущего API Sprint 2
    description TEXT NOT NULL,
    image_url TEXT NOT NULL,
    category_id INT NOT NULL REFERENCES categories(id) ON DELETE RESTRICT
);

-- CARTS
CREATE TABLE IF NOT EXISTS carts (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    amount INT NOT NULL DEFAULT 0,
    amount_cost NUMERIC(12,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- CART_PRODUCTS (join)
CREATE TABLE IF NOT EXISTS cart_products (
    id SERIAL PRIMARY KEY,
    cart_id INT NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
    product_id INT NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    quantity INT NOT NULL CHECK (quantity > 0) DEFAULT 1,
    UNIQUE (cart_id, product_id)
);

-- SEED Sprint 2 (категории + продукты)
INSERT INTO categories (name) VALUES
('Корм для собак'),
('Корм для котів'),
('Рослини для дому')
ON CONFLICT (name) DO NOTHING;

INSERT INTO products (name, cost, producer, type, guarantee, weight, unit, title, description, image_url, category_id)
SELECT
	'Сухий корм преміум' AS name,
	799.00 AS cost,
	'Savory' AS producer,
	c.name AS type,
	NULL AS guarantee,
	NULL AS weight,
	NULL AS unit,
	'Сухий корм преміум' AS title,
	'Збалансований корм для дорослих собак' AS description,
	'https://example.com/dog-food.jpg' AS image_url,
	c.id AS category_id
FROM categories c
WHERE c.name = 'Корм для собак'
  AND NOT EXISTS (
    SELECT 1 FROM products p WHERE p.name = 'Сухий корм преміум' AND p.category_id = c.id
  );

INSERT INTO products (name, cost, producer, type, guarantee, weight, unit, title, description, image_url, category_id)
SELECT
	'Вологий корм для котів',
	499.00,
	'Savory',
	c.name,
	NULL,
	NULL,
	NULL,
	'Вологий корм для котів',
	'Мʼясний вологий корм у паучах',
	'https://example.com/cat-food.jpg',
	c.id
FROM categories c
WHERE c.name = 'Корм для котів'
  AND NOT EXISTS (
    SELECT 1 FROM products p WHERE p.name = 'Вологий корм для котів' AND p.category_id = c.id
  );

INSERT INTO products (name, cost, producer, type, guarantee, weight, unit, title, description, image_url, category_id)
SELECT
	'Монстера Deliciosa',
	1299.00,
	'Green',
	c.name,
	NULL,
	NULL,
	NULL,
	'Монстера Deliciosa',
	'Декоративна кімнатна рослина середнього розміру',
	'https://example.com/monstera.jpg',
	c.id
FROM categories c
WHERE c.name = 'Рослини для дому'
  AND NOT EXISTS (
    SELECT 1 FROM products p WHERE p.name = 'Монстера Deliciosa' AND p.category_id = c.id
  );

-- Обнуляем кошики после пересборки.
DELETE FROM carts;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cart_products;
DROP TABLE IF EXISTS carts;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS users;
-- categories не трогаем, чтобы не ломать ранние спринты
-- +goose StatementEnd
