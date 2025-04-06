# Тестовое задание Junior Golang Developer в Effective Mobile

## Накатить миграции

`goose -dir migrations/ postgres up`

## Генерация swagger и запуск сервера

`swag init && go run main.go`

- [x] Выставлены rest методы
- [x] Данные обогащаются из стороннего API
- [x] Данные хранятся в pg (схема БД cоздана путем миграций goose)
- [x] Код покрыт debug- и info-логами (автоматом от Gin)
- [x] Конфигурационные данные вынесены в .env
- [x] Сгенерирован сваггер на реализованное API (swaggo/swag)

### Путь к swagger

http://localhost:8080/swagger/index.html
