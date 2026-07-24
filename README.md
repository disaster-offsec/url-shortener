# URL Shortener Service

Микросервис для сокращения ссылок, реализованный на Go.

## Особенности
- In-memory хранилище с двумя `map` для O(1) поиска в обе стороны
- Потокобезопасность через `sync.RWMutex`
- Проверка существующих URL — одна ссылка = один короткий код
- Криптостойкая генерация коротких кодов (`crypto/rand`)
- Контейнеризация через Docker
- Обработка ошибок (400, 404)
- Идемпотентность — повторные запросы возвращают ту же ссылку

## Запуск

### Через Docker (рекомендуется)

```bash
docker compose up --build
```
Сервер будет доступен по адресу: http://localhost:8080

### Локально

```bash
go mod tidy
go run main.go storage.go
```

# API

POST /shorten — создать короткую ссылку

Запрос:

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com"}'
```
Ответ:

```json
{"short_url":"http://localhost:8080/DhH9xw"}
```

Если ссылка уже существует, вернётся та же короткая ссылка.

GET /{code} — редирект по короткой ссылке

```bash
curl -v http://localhost:8080/DhH9xw
```

Ожидаемый ответ: HTTP 302 Found с заголовком Location: https://example.com

Ошибки
| Сценарий | HTTP статус |
| Невалидный JSON |	400 Bad Request |
| Поле url пустое |	400 Bad Request |
| Короткий код не найден | 404 Not Found |

# Архитектура

Хранилище
Используются две map для обеспечения двунаправленного поиска за O(1):

`codeToURL` — для редиректа по короткому коду

`urlToCode` — для проверки существования URL

Все операции потокобезопасны благодаря sync.RWMutex.

Генерация коротких кодов
Код генерируется через crypto/rand и кодируется в base64.URLEncoding. Длина — 6 символов. Этого достаточно для ~56 миллиардов комбинаций.

# Технологии

Go 1.22

Стандартная библиотека (net/http, sync, crypto/rand)

Docker + docker-compose

## Структура проекта

```
.
├── main.go          # HTTP-сервер и хэндлеры
├── storage.go       # In-memory хранилище
├── Dockerfile       # Многоступенчатая сборка
├── docker-compose.yml
├── go.mod
└── README.md
```
