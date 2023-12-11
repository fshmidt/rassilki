# Приложение для управления рассылками

<!-- ToC start -->
# Содержание

1. [Описание задачи](#Описание-задачи)
2. [Реализация](#Реализация)
3. [Endpoints](#Endpoints)
4. [Запуск](#Запуск)
5. [Swagger UI](#Swagger UI)
<!-- ToC end -->

# Описание задачи

Разработать API для управления рассылками и клиентами. Приложение предоставляет функциональность по созданию, редактированию и удалению рассылок, а также добавлению, редактированию и удалению клиентов и осуществляет рассылки пользуясь внешним API.

# Реализация

- Применение принципов REST API.
- Применение принципов "Чистой Архитектуры" и техники внедрения зависимостей.
- Использование фреймворка [gin-gonic/gin](https://github.com/gin-gonic/gin) для обработки HTTP-запросов.
- Взаимодействие с СУБД Postgres с использованием библиотеки [sqlx](https://github.com/jmoiron/sqlx) и написание SQL-запросов.
- Конфигурация приложения осуществляется с использованием библиотеки [viper](https://github.com/spf13/viper).
- Запуск базы данных в Docker-контейнере.
- Асинхронная работа рассылок.
- Swagger UI

**Структура проекта:**

```
.
├── pkg
│   ├── handler     // обработчики запросов
│   ├── service     // бизнес-логика
│   └── repository  // взаимодействие с БД
├── cmd             // точка входа в приложение
├── server          // серверная часть и логика обращений к внешнему API
├── schema          // SQL файлы с миграциями
├── configs         // файлы конфигурации
```


# Endpoints

## Клиенты

- **Создание клиента**
  - Endpoint: `POST /clients`
  - Тело запроса:
    ```json
    {
        "phone": "7123456",
        "code": "123",
        "tag": "important",
        "timezone": "UTC"
    }
    ```
  - Успешный ответ:
    ```json
    {"id": 1}
    ```

- **Получение клиента по ID**
  - Endpoint: `GET /clients/:id`
  - Успешный ответ:
    ```json
    {
        "phone": "7123456",
        "code": "123",
        "tag": "important",
        "timezone": "UTC"
    }
    ```

- **Обновление клиента по ID**
  - Endpoint: `PUT /clients/:id`
  - Тело запроса:
    ```json
    {
        "phone": "7666666",
        "tag": "urgent"
    }
    ```
  - Успешный ответ:
    ```json
    {"status": "ok"}
    ```

- **Удаление клиента по ID**
  - Endpoint: `DELETE /clients/:id`
  - Успешный ответ:
    ```json
    {"status": "ok"}
    ```

## Рассылки

- **Создание рассылки**
  - Endpoint: `POST /rassilki`
  - Тело запроса:
    ```json
    {
        "start-time": "2023-12-31T23:59:59Z",
        "message": "Special offer for our clients!",
        "filter": ["important", "rich"],
        "end-time": "2024-01-15T23:59:59Z"
    }
    ```
  - Успешный ответ:
    ```json
    {"id": 1}
    ```

- **Получение статистики рассылки по ID**
  - Endpoint: `GET /rassilki/:id`
  - Успешный ответ:
    ```json
    {
    "ras_id": 7,
    "total": 20,
    "sent": 20,
    "not_sent": 0
    }
    ```

- **Получение статистики рпо всем рассылкам**
  - Endpoint: `GET /rassilki/`
  - Успешный ответ:
    ```
    [
        {
            "ras_id": 3,
            "total": 15,
            "sent": 15,
            "not_sent": 0
        },
        {
            "ras_id": 7,
            "total": 20,
            "sent": 20,
            "not_sent": 0
        },
        {
            "ras_id": 8,
            "total": 20,
            "sent": 0,
            "not_sent": 20
        }
    ]
    ```

- **Обновление рассылки по ID**
  - Endpoint: `PUT /rassilki/:id`
  - Тело запроса:
    ```json
    {
        "start-time": "2024-01-01T00:00:00Z",
        "message": "Updated special offer!",
        "filter": ["important", "rich","stupid"],
        "end-time": "2024-01-20T23:59:59Z"
    }
    ```
  - Успешный ответ:
    ```json
    {"status": "ok"}
    ```

- **Удаление рассылки по ID**
  - Endpoint: `DELETE /rassilki/:id`
  - Успешный ответ:
    ```json
    {"status": "ok"}
    ```

# Запуск

```
#bash
go run cmd/main.go

```

Если приложение запускается впервые, необходимо применить миграции к базе данных:

```
docker pull postgres
docker run --name=rassilki-db -e POSTGRES_PASSWORD='qwerty' -p 5436:5432 -d --rm postgres
migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' up
go mod tidy

```
Так же необходимо иметь переменные окружения JWT с JWT-токеном доступа и DB_PASSWORD с паролем от базы данных (для примера и как видно выше это `qwerty`)

# Swagger UI

По адресу /docs/index.html находится страница со Swagger Ui:
![Swagger Ui](./docs/swagger.png)
