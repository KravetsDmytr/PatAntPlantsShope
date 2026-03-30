-- +goose Up
-- +goose StatementBegin
INSERT INTO categories (name) VALUES
('Корм для собак'),
('Корм для котів'),
('Рослини для дому')
ON CONFLICT (name) DO NOTHING;

INSERT INTO products (name, description, price, image_url, category_id)
SELECT
	'Сухий корм преміум',
	'Збалансований корм для дорослих собак',
	799.00,
	'https://example.com/dog-food.jpg',
	c.id
FROM categories c
WHERE c.name = 'Корм для собак'
  AND NOT EXISTS (
    SELECT 1 FROM products p WHERE p.name = 'Сухий корм преміум' AND p.category_id = c.id
  );

INSERT INTO products (name, description, price, image_url, category_id)
SELECT
	'Вологий корм для котів',
	'Мʼясний вологий корм у паучах',
	499.00,
	'https://example.com/cat-food.jpg',
	c.id
FROM categories c
WHERE c.name = 'Корм для котів'
  AND NOT EXISTS (
    SELECT 1 FROM products p WHERE p.name = 'Вологий корм для котів' AND p.category_id = c.id
  );

INSERT INTO products (name, description, price, image_url, category_id)
SELECT
	'Монстера Deliciosa',
	'Декоративна кімнатна рослина середнього розміру',
	1299.00,
	'https://example.com/monstera.jpg',
	c.id
FROM categories c
WHERE c.name = 'Рослини для дому'
  AND NOT EXISTS (
    SELECT 1 FROM products p WHERE p.name = 'Монстера Deliciosa' AND p.category_id = c.id
  );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM products WHERE image_url IN (
	'https://example.com/dog-food.jpg',
	'https://example.com/cat-food.jpg',
	'https://example.com/monstera.jpg'
);
DELETE FROM categories WHERE name IN ('Корм для собак', 'Корм для котів', 'Рослини для дому');
-- +goose StatementEnd
