package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"website-dm/internal/auth"
	"website-dm/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	CreateUser(ctx context.Context, req models.RegisterRequest, passHash string) (int, error)
	GetUserCredentials(ctx context.Context, login string) (int, string, error)
	ListCategories(ctx context.Context) ([]models.Category, error)
	ListProducts(ctx context.Context, f models.ProductFilters) ([]models.Product, error)
	GetProductByID(ctx context.Context, id int) (*models.Product, error)
	GetOrCreateCartID(ctx context.Context, userID int) (int, error)
	UpsertCartItem(ctx context.Context, cartID, productID, qty int) error
	GetCart(ctx context.Context, userID int) (*models.CartResponse, error)
}

type Service struct {
	repo      Repository
	jwtSecret string
}

func New(repo Repository, jwtSecret string) *Service {
	return &Service{repo: repo, jwtSecret: jwtSecret}
}

func (s *Service) Register(ctx context.Context, req models.RegisterRequest) error {
	if strings.TrimSpace(req.Login) == "" || strings.TrimSpace(req.Password) == "" || strings.TrimSpace(req.Email) == "" {
		return errors.New("login, email та password є обов'язковими")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = s.repo.CreateUser(ctx, req, string(hash))
	return err
}

func (s *Service) Login(ctx context.Context, req models.LoginRequest) (string, error) {
	id, hash, err := s.repo.GetUserCredentials(ctx, req.Login)
	if err != nil {
		return "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
		return "", errors.New("невірний логін або пароль")
	}

	return auth.GenerateToken(s.jwtSecret, id)
}

func (s *Service) ValidateToken(token string) (int, error) {
	return auth.ValidateToken(s.jwtSecret, token)
}

func (s *Service) GetCategories(ctx context.Context) ([]models.Category, error) {
	return s.repo.ListCategories(ctx)
}

func (s *Service) GetProducts(ctx context.Context, f models.ProductFilters) ([]models.Product, error) {
	return s.repo.ListProducts(ctx, f)
}

func (s *Service) GetProductByID(ctx context.Context, id int) (*models.Product, error) {
	return s.repo.GetProductByID(ctx, id)
}

func (s *Service) AddToCart(ctx context.Context, userID int, req models.AddToCartRequest) error {
	if req.Quantity <= 0 {
		return errors.New("quantity має бути більше 0")
	}
	if _, err := s.repo.GetProductByID(ctx, req.ProductID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("товар не знайдено")
		}
		return err
	}

	cartID, err := s.repo.GetOrCreateCartID(ctx, userID)
	if err != nil {
		return err
	}
	return s.repo.UpsertCartItem(ctx, cartID, req.ProductID, req.Quantity)
}

func (s *Service) GetCart(ctx context.Context, userID int) (*models.CartResponse, error) {
	return s.repo.GetCart(ctx, userID)
}
