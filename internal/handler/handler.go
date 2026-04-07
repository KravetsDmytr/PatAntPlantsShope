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

// Register registers a new user
// @Summary Реєстрація користувача
// @Description Створює нового користувача в системі
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Дані користувача"
// @Success 201 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /auth/register [post]
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

// Login authenticates user and returns JWT token
// @Summary Авторизація користувача
// @Description Виконує вхід користувача та повертає JWT токен
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Облікові дані"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Router /auth/login [post]
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

// Categories returns all categories
// @Summary Отримати категорії
// @Description Повертає список категорій товарів
// @Tags Categories
// @Produce json
// @Success 200 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /categories [get]
func (h *Handler) Categories(c *gin.Context) {
	data, err := h.service.GetCategories(c.Request.Context())
	if err != nil {
		respond(c, http.StatusInternalServerError, "Не вдалося отримати категорії", nil, err.Error())
		return
	}
	respond(c, http.StatusOK, "Категорії отримано", data, nil)
}

// Products returns products list with filters
// @Summary Отримати список товарів
// @Description Повертає товари з фільтрацією за категорією, вартістю та пошуковим запитом
// @Tags Products
// @Produce json
// @Param category_id query int false "ID категорії"
// @Param min_cost query number false "Мінімальна вартість"
// @Param max_cost query number false "Максимальна вартість"
// @Param q query string false "Пошуковий запит"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /products [get]
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

// ProductByID returns product by id
// @Summary Отримати товар за ID
// @Description Повертає детальну інформацію про товар
// @Tags Products
// @Produce json
// @Param id path int true "ID товару"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /products/{id} [get]
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

// AddToCart adds product to cart
// @Summary Додати товар у кошик
// @Description Додає товар у кошик авторизованого користувача
// @Tags Cart
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT токен"
// @Param request body models.AddToCartRequest true "Товар і кількість"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Router /cart/items [post]
// @Security ApiKeyAuth
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

// Cart returns user cart
// @Summary Отримати кошик
// @Description Повертає вміст кошика авторизованого користувача
// @Tags Cart
// @Produce json
// @Param Authorization header string true "Bearer JWT токен"
// @Success 200 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /cart [get]
// @Security ApiKeyAuth
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
