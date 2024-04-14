# Тестовое задание для стажёра Backend

<!-- ToC start -->
# Описание задачи
https://github.com/avito-tech/backend-trainee-assignment-2024?tab=readme-ov-file

# Реализация
- Применение фреймворка [gin-gonic/gin](https://github.com/gin-gonic/gin).
- Применение СУБД Postgres посредствов библиотеки [sqlx](https://github.com/jmoiron/sqlx) и написанием SQL запросов.
- Применение СУБД Redis посредством библиотеки [go-redis](https://github.com/redis/go-redis/)
- Контейнеризация с помощью Docker и docker-compose

**Структура проекта:**
```
.
├── internal
│   ├── app         // точка запуска приложения
│   ├── cache       // настройка подключения к кэшу
│   ├── config      // общие конфигурации приложения
│   ├── delivery    // слой обработки запросов
│   ├── models      // структуры сущностей приложения
│   ├── service     // слой бизнес-логики
│   └── repository  // слой взаимодействия с БД
├── cmd             // точка входа в приложение
└── db              // SQL файлы миграции
```

# Endpoints
Реализованы предоставленные [api](https://drive.google.com/file/d/1l4PMTPzsjksRCd_lIm0mVfh4U0Jn-A2R/view)
Для дополнительных заданий также реализованы:
- DELETE /api/banner
  - in query:
    - tag_id - идентификатор тега
    - feature_id - идентификатор фичи
- GET /api/banner/:id
  - in path:
    - id - идентификатор баннера
- PUT /api/banner/:id
  - in path:
    - id - идентификатор баннера
  - in request body
    - version - номер версии баннера
   
🟡 Также для управления токенами реализована авторизация
- POST /auth/sigh-in
  - in request body:
    - login - логин
    - password - пароль
    
# Запуск
```
make build && make run
```
Если приложение запускается впервые, необходимо применить миграции к базе данных:
```
make migrate-up
```
🟡 Также мигрированы логины и пароли пользователей и админа для опробывания сервиса (их можно найти в файлах мигрирования)

# Примеры
Запросы сгенерированы командой curl
### 1. POST /auth/sign-in/
**Запрос:**
```
$ curl --location --request POST 'localhost:8000/auth/sign-in' \
--header 'Content-Type: application/json' \
--data-raw '{
    "login": "admin",
    "password": "admin!"
}'
```
**Тело ответа:**
```
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJNYXBDbGFpbXMiOnsiZXhwIjoxNzExNjA0OTY0LCJpYXQiOjE3MTE1NjE3NjR9LCJsb2dpbiI6ImFsaWJlayJ9.QsIocl01gHVaUuTtQEhOJgmm6Mgu-K0LwMddYPKT7v4"
}
```
### 2. POST /api/banner/
**Запрос:**
```
$ curl --location --request POST 'localhost:8000/api/banner' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJNYXBDbGFpbXMiOnsiZXhwIjoxNzExNjA0OTY0LCJpYXQiOjE3MTE1NjE3NjR9LCJsb2dpbiI6ImFsaWJlayJ9.QsIocl01gHVaUuTtQEhOJgmm6Mgu-K0LwMddYPKT7v4' \
--data-raw '{
   "tag_ids": [1, 2, 3],
   "feature_id": 1,
   "content": {
       "title": "some",
       "text": "some",
       "url": "some"
   },
   "is_active": true
}'
```
**Тело ответа:**
```
{
    "banner_id": 2
}
```
### 3. GET /api/banner/
**Запрос:**
```
$ curl --location --request GET 'localhost:8000/api/banner?tag_id=7&limit=1' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJNYXBDbGFpbXMiOnsiZXhwIjoxNzExNjA0OTY0LCJpYXQiOjE3MTE1NjE3NjR9LCJsb2dpbiI6ImFsaWJlayJ9.QsIocl01gHVaUuTtQEhOJgmm6Mgu-K0LwMddYPKT7v4'
```
**Тело ответа:**
```
[
    {
        "id": 1,
        "content": {
            "title": "sss",
            "text": "sss",
            "url": "sss"
        },
        "tag_ids": [
            7,
            8,
            5
        ],
        "feature_id": 2,
        "is_active": true,
        "created_at": "2024-04-14T17:42:14.262925Z",
        "updated_at": "2024-04-14T17:45:00.318351Z"
    }
]
```
### 4. PATCH /api/banner/:id
**Запрос:**
```
$ curl --location --request PATCH 'localhost:8000/api/banner/2' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJNYXBDbGFpbXMiOnsiZXhwIjoxNzExNjA0OTY0LCJpYXQiOjE3MTE1NjE3NjR9LCJsb2dpbiI6ImFsaWJlayJ9.QsIocl01gHVaUuTtQEhOJgmm6Mgu-K0LwMddYPKT7v4' \
--data-raw '{
    "tag_ids": [7, 8, 5],
   "feature_id": 3,
    "content": {
       "title": "updated",
       "text": "updated",
       "url": "updated"
   }
}'
```
**Тело ответа:**
```
{
    "description": "ok"
}
```
### 5. DELETE /api/banner/:id
**Запрос:**
```
$ curl --location --request DELETE 'localhost:8000/api/banner/2' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJNYXBDbGFpbXMiOnsiZXhwIjoxNzExNjA0OTY0LCJpYXQiOjE3MTE1NjE3NjR9LCJsb2dpbiI6ImFsaWJlayJ9.QsIocl01gHVaUuTtQEhOJgmm6Mgu-K0LwMddYPKT7v4''
```
**Тело ответа:**
```
{
    "description": "Баннер успешно удален"
}
```
### 6. DELETE /api/banner/
**Запрос:**
```
$ curl --location --request DELETE 'localhost:8000/api/banner?feature_id=2' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJNYXBDbGFpbXMiOnsiZXhwIjoxNzExNjA0OTY0LCJpYXQiOjE3MTE1NjE3NjR9LCJsb2dpbiI6ImFsaWJlayJ9.QsIocl01gHVaUuTtQEhOJgmm6Mgu-K0LwMddYPKT7v4''
```
**Тело ответа:**
```
{
    "description": "Баннеры успешно удален"
}
```
### 7. GET /api/banner/:id
**Запрос:**
```
$ curl --location --request GET 'localhost:8000/api/banner/3' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJNYXBDbGFpbXMiOnsiZXhwIjoxNzExNjA0OTY0LCJpYXQiOjE3MTE1NjE3NjR9LCJsb2dpbiI6ImFsaWJlayJ9.QsIocl01gHVaUuTtQEhOJgmm6Mgu-K0LwMddYPKT7v4''
```
**Тело ответа:**
```
[
    {
        "id": 3,
        "content": {
            "title": "first",
            "text": "first",
            "url": "first"
        },
        "tag_ids": [
            1,
            2,
            3
        ],
        "feature_id": 1,
        "is_active": true,
        "created_at": "2024-04-14T19:28:09.176156Z",
        "updated_at": "2024-04-14T19:28:09.176156Z"
    },
    {
        "id": 3,
        "content": {
            "title": "second",
            "text": "second",
            "url": "second"
        },
        "tag_ids": [
            4,
            5
        ],
        "feature_id": 4,
        "is_active": true,
        "created_at": "2024-04-14T19:28:09.176156Z",
        "updated_at": "2024-04-14T19:28:54.150835Z"
    }
]
```
### 8. PUT /api/banner/:id
**Запрос:**
```
$ curl --location --request PUT 'localhost:8000/api/banner/3' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJNYXBDbGFpbXMiOnsiZXhwIjoxNzExNjA0OTY0LCJpYXQiOjE3MTE1NjE3NjR9LCJsb2dpbiI6ImFsaWJlayJ9.QsIocl01gHVaUuTtQEhOJgmm6Mgu-K0LwMddYPKT7v4' \
--data-raw '{
    "version": 1
}'
```
**Тело ответа:**
```
{
    "description": "ok"
}
```
### 9. PUT /api/user_banner
**Запрос:**
```
$ curl --location --request PUT 'localhost:8000/api/user_banner?tag_id=1&feature_id=1' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJNYXBDbGFpbXMiOnsiZXhwIjoxNzExNjA0OTY0LCJpYXQiOjE3MTE1NjE3NjR9LCJsb2dpbiI6ImFsaWJlayJ9.QsIocl01gHVaUuTtQEhOJgmm6Mgu-K0LwMddYPKT7v4' \
--data-raw '{
    "version": 1
}'
```
**Тело ответа:**
```
{
    "title": "first",
    "text": "first",
    "url": "first"
}
```

# Проблема
🔴 Не успел написать ни одного теста
