# Quotes App

## Мини-сервис “Цитатник”


## Запуск проекта

1. Клонируйте репозиторий:

```bash
git clone https://github.com/Shammaaaa/Quotes.git
cd Quotes
```


2. Установите зависимости:

```bash
go mod download
```
3. Запуск сервера:


```bash
go run cmd/server/main.go
```
4. Запуск тестов

```bash
go test -v ./...
```

После запуска приложение будет доступно на http://localhost:8080


## API Endpoints
Создание цитаты

Метод: POST

URL: /quotes

Тело запроса:
```json
{
  "author": "Confucius",
  "quote": "Life is really simple, but we insist on making it complicated."
}
```
Ответ:
```json
{
  "id": 1,
  "author": "Confucius",
  "quote": "Life is really simple, but we insist on making it complicated.",
  "created_at": "2025-05-24T11:24:27.2535694+03:00"
}
```
Получение цитат

Метод: GET

URL: /quotes

Ответ:
```json
{
  "id": 1,
  "author": "Confucius",
  "quote": "Life is really simple, but we insist on making it complicated.",
  "created_at": "2025-05-24T11:24:27.2535694+03:00"
}
```
Получение случайной цитаты

Метод: GET

URL: /quotes/random

Ответ:
```json
{
  "id": 1,
  "author": "Confucius",
  "quote": "Life is really simple, but we insist on making it complicated.",
  "created_at": "2025-05-24T11:24:27.2535694+03:00"
}
```
Получение цитаты по автору

URL: /quotes?author=Confucius

Метод: GET


Ответ: 
```json
{
  "id": 3,
  "author": "Confucius",
  "quote": "Life is really simple, but we insist on making it complicated.",
  "created_at": "2025-05-24T11:24:27.2535694+03:00"
}
```

Удаление цитаты

URL: /quotes/1

Метод: DELETE

Ответ: 204 No Content

