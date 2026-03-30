package handler

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"website-dm/internal/models"

	"github.com/gin-gonic/gin"
)

type Service interface {
	Register(ctx context.Context, req models.RegisterRequest) error
	Login(ctx context.Context, req models.LoginRequest) (string, error)
	GetCategories(ctx context.Context) ([]models.Category, error)
	GetProducts(ctx context.Context, f models.ProductFilters) ([]models.Product, error)
	GetProductByID(ctx context.Context, id int) (*models.Product, error)
	AddToCart(ctx context.Context, userID int, req models.AddToCartRequest) error
	GetCart(ctx context.Context, userID int) (*models.CartResponse, error)
}

type Handler struct {
	service Service
}

func New(service Service) *Handler {
	return &Handler{service: service}
}

func respond(c *gin.Context, code int, message string, data interface{}, err interface{}) {
	c.JSON(code, models.APIResponse{Message: message, Data: data, Error: err})
}

func (h *Handler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond(c, http.StatusBadRequest, "Невірний запит", nil, err.Error())
		return
	}
	if err := h.service.Register(c.Request.Context(), req); err != nil {
		respond(c, http.StatusBadRequest, "Не вдалося створити користувача", nil, err.Error())
		return
	}
	respond(c, http.StatusCreated, "Користувача створено", nil, nil)
}

func (h *Handler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond(c, http.StatusBadRequest, "Невірний запит", nil, err.Error())
		return
	}
	token, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		respond(c, http.StatusUnauthorized, "Помилка авторизації", nil, err.Error())
		return
	}
	respond(c, http.StatusOK, "Вхід успішний", models.TokenResponse{Token: token}, nil)
}

func (h *Handler) Categories(c *gin.Context) {
	data, err := h.service.GetCategories(c.Request.Context())
	if err != nil {
		respond(c, http.StatusInternalServerError, "Не вдалося отримати категорії", nil, err.Error())
		return
	}
	respond(c, http.StatusOK, "Категорії отримано", data, nil)
}

func (h *Handler) Products(c *gin.Context) {
	var filters models.ProductFilters

	if raw := c.Query("category_id"); raw != "" {
		id, err := strconv.Atoi(raw)
		if err != nil {
			respond(c, http.StatusBadRequest, "Невірний category_id", nil, err.Error())
			return
		}
		filters.CategoryID = &id
	}

	// ERD/БД: cost. Для удобства принимаем и старые параметры min_price/max_price.
	if raw := c.Query("min_cost"); raw != "" {
		v, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			respond(c, http.StatusBadRequest, "Невірний min_cost", nil, err.Error())
			return
		}
		filters.MinCost = &v
	} else if raw := c.Query("min_price"); raw != "" {
		v, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			respond(c, http.StatusBadRequest, "Невірний min_price", nil, err.Error())
			return
		}
		filters.MinCost = &v
	}

	if raw := c.Query("max_cost"); raw != "" {
		v, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			respond(c, http.StatusBadRequest, "Невірний max_cost", nil, err.Error())
			return
		}
		filters.MaxCost = &v
	} else if raw := c.Query("max_price"); raw != "" {
		v, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			respond(c, http.StatusBadRequest, "Невірний max_price", nil, err.Error())
			return
		}
		filters.MaxCost = &v
	}

	filters.Query = c.Query("q")

	data, err := h.service.GetProducts(c.Request.Context(), filters)
	if err != nil {
		respond(c, http.StatusInternalServerError, "Не вдалося отримати товари", nil, err.Error())
		return
	}
	respond(c, http.StatusOK, "Товари отримано", data, nil)
}

func (h *Handler) ProductByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respond(c, http.StatusBadRequest, "Невірний id", nil, err.Error())
		return
	}
	data, err := h.service.GetProductByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respond(c, http.StatusNotFound, "Товар не знайдено", nil, "not found")
			return
		}
		respond(c, http.StatusInternalServerError, "Не вдалося отримати товар", nil, err.Error())
		return
	}
	respond(c, http.StatusOK, "Товар отримано", data, nil)
}

func (h *Handler) AddToCart(c *gin.Context) {
	userIDAny, _ := c.Get("user_id")
	userID, ok := userIDAny.(int)
	if !ok {
		respond(c, http.StatusUnauthorized, "Потрібна авторизація", nil, "invalid user")
		return
	}

	var req models.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond(c, http.StatusBadRequest, "Невірний запит", nil, err.Error())
		return
	}

	if err := h.service.AddToCart(c.Request.Context(), userID, req); err != nil {
		respond(c, http.StatusBadRequest, "Не вдалося додати товар до кошика", nil, err.Error())
		return
	}
	respond(c, http.StatusOK, "Товар додано до кошика", nil, nil)
}

func (h *Handler) Cart(c *gin.Context) {
	userIDAny, _ := c.Get("user_id")
	userID, ok := userIDAny.(int)
	if !ok {
		respond(c, http.StatusUnauthorized, "Потрібна авторизація", nil, "invalid user")
		return
	}

	data, err := h.service.GetCart(c.Request.Context(), userID)
	if err != nil {
		respond(c, http.StatusInternalServerError, "Не вдалося отримати кошик", nil, err.Error())
		return
	}
	respond(c, http.StatusOK, "Кошик отримано", data, nil)
}
