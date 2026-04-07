package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"website-dm/internal/models"
)

type Repository struct {
	db *sql.DB
}

type productScanner interface {
	Scan(dest ...interface{}) error
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func scanProduct(s productScanner) (models.Product, error) {
	var p models.Product
	var title sql.NullString
	var weight sql.NullFloat64
	var unit sql.NullString
	var guarantee sql.NullTime

	err := s.Scan(
		&p.ID,
		&p.Name,
		&p.Producer,
		&p.Type,
		&p.Description,
		&title,
		&p.Cost,
		&weight,
		&unit,
		&guarantee,
		&p.ImageURL,
		&p.CategoryID,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return models.Product{}, err
	}

	if title.Valid {
		p.Title = &title.String
	}
	if weight.Valid {
		p.Weight = &weight.Float64
	}
	if unit.Valid {
		p.Unit = &unit.String
	}
	if guarantee.Valid {
		p.Guarantee = &guarantee.Time
	}

	return p, nil
}

func (r *Repository) CreateUser(ctx context.Context, req models.RegisterRequest, passHash string) (int, error) {
	const q = `
		INSERT INTO users (login, first_name, last_name, email, password_hash)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`
	var id int
	err := r.db.QueryRowContext(ctx, q, req.Login, req.FirstName, req.LastName, req.Email, passHash).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repository) GetUserCredentials(ctx context.Context, login string) (int, string, error) {
	const q = `SELECT id, password_hash FROM users WHERE login = $1`
	var id int
	var hash string
	err := r.db.QueryRowContext(ctx, q, login).Scan(&id, &hash)
	if err != nil {
		return 0, "", err
	}
	return id, hash, nil
}

func (r *Repository) ListCategories(ctx context.Context) ([]models.Category, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name FROM categories ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]models.Category, 0)
	for rows.Next() {
		var c models.Category
		if err = rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (r *Repository) ListProducts(ctx context.Context, f models.ProductFilters) ([]models.Product, error) {
	base := `SELECT id, name, producer, type, description, title, cost::double precision AS cost,
	                weight::double precision, unit, guarantee, image_url, category_id, created_at, updated_at
	         FROM products`
	where := make([]string, 0)
	args := make([]interface{}, 0)

	if f.CategoryID != nil {
		args = append(args, *f.CategoryID)
		where = append(where, fmt.Sprintf("category_id = $%d", len(args)))
	}
	if f.MinCost != nil {
		args = append(args, *f.MinCost)
		where = append(where, fmt.Sprintf("cost >= $%d", len(args)))
	}
	if f.MaxCost != nil {
		args = append(args, *f.MaxCost)
		where = append(where, fmt.Sprintf("cost <= $%d", len(args)))
	}
	if strings.TrimSpace(f.Query) != "" {
		args = append(args, "%"+strings.TrimSpace(f.Query)+"%")
		where = append(where, fmt.Sprintf("name ILIKE $%d", len(args)))
	}

	query := base
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}
	query += " ORDER BY id"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		p, scanErr := scanProduct(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

func (r *Repository) GetProductByID(ctx context.Context, id int) (*models.Product, error) {
	const q = `SELECT id, name, producer, type, description, title, cost::double precision AS cost,
	                  weight::double precision, unit, guarantee, image_url, category_id, created_at, updated_at
	           FROM products WHERE id = $1`
	p, err := scanProduct(r.db.QueryRowContext(ctx, q, id))
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repository) GetOrCreateCartID(ctx context.Context, userID int) (int, error) {
	const getQ = `SELECT id FROM carts WHERE user_id = $1`
	var cartID int
	err := r.db.QueryRowContext(ctx, getQ, userID).Scan(&cartID)
	if err == nil {
		return cartID, nil
	}
	if err != sql.ErrNoRows {
		return 0, err
	}

	const createQ = `INSERT INTO carts (user_id) VALUES ($1) RETURNING id`
	if err = r.db.QueryRowContext(ctx, createQ, userID).Scan(&cartID); err != nil {
		return 0, err
	}
	return cartID, nil
}

func (r *Repository) UpsertCartItem(ctx context.Context, cartID, productID, qty int) error {
	const q = `
		INSERT INTO cart_products (cart_id, product_id, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT (cart_id, product_id)
		DO UPDATE SET quantity = cart_products.quantity + EXCLUDED.quantity`

	if _, err := r.db.ExecContext(ctx, q, cartID, productID, qty); err != nil {
		return err
	}

	const updateAgg = `
		UPDATE carts c
		SET
			amount = agg.qty,
			amount_cost = agg.total_cost,
			updated_at = NOW()
		FROM (
			SELECT
				cp.cart_id,
				SUM(cp.quantity) AS qty,
				SUM(p.cost::double precision * cp.quantity::double precision) AS total_cost
			FROM cart_products cp
			JOIN products p ON p.id = cp.product_id
			WHERE cp.cart_id = $1
			GROUP BY cp.cart_id
		) agg
		WHERE c.id = agg.cart_id`

	_, err := r.db.ExecContext(ctx, updateAgg, cartID)
	return err
}

func (r *Repository) GetCart(ctx context.Context, userID int) (*models.CartResponse, error) {
	const q = `
		SELECT cp.product_id,
		       p.name,
		       p.cost::double precision AS cost,
		       cp.quantity,
		       (p.cost::double precision * cp.quantity::double precision) AS line_total
		FROM carts c
		JOIN cart_products cp ON cp.cart_id = c.id
		JOIN products p ON p.id = cp.product_id
		WHERE c.user_id = $1
		ORDER BY cp.id`

	rows, err := r.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resp := &models.CartResponse{Items: make([]models.CartItem, 0)}
	for rows.Next() {
		var item models.CartItem
		if err = rows.Scan(&item.ProductID, &item.ProductName, &item.Cost, &item.Quantity, &item.LineTotal); err != nil {
			return nil, err
		}
		resp.AmountCost += item.LineTotal
		resp.Items = append(resp.Items, item)
	}
	return resp, rows.Err()
}
