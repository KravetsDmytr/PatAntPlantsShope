package models

import "time"

type APIResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   interface{} `json:"error"`
}

type RegisterRequest struct {
	Login     string `json:"login"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Cost        float64   `json:"cost"`
	ImageURL    string    `json:"image_url"`
	CategoryID  int       `json:"category_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type ProductFilters struct {
	CategoryID *int
	MinCost    *float64
	MaxCost    *float64
	Query      string
}

type AddToCartRequest struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type CartItem struct {
	ProductID   int     `json:"product_id"`
	ProductName string  `json:"product_name"`
	Cost        float64 `json:"cost"`
	Quantity    int     `json:"quantity"`
	LineTotal   float64 `json:"line_total"`
}

type CartResponse struct {
	Items []CartItem `json:"items"`
	AmountCost float64 `json:"amount_cost"`
}
