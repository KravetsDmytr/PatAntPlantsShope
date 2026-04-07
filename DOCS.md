# Автоматизована генерація документації (ЛР №7)

У проєкті підключено Swagger UI через `swaggo`:
- пакет документації: `api/openapi`
- runtime-налаштування через `docs.SwaggerInfo.*` в `internal/app/app.go`

## Відкрити документацію

Після запуску API:

- `http://localhost:8080/swagger/index.html`

## Оновити документацію з коду (автогенерація)

1. Встановити CLI:

`go install github.com/swaggo/swag/cmd/swag@latest`

2. Згенерувати docs на основі анотацій:

`swag init -g cmd/main.go --output ./api/openapi`

Після цього swagger-спека буде оновлена автоматично.
