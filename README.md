```markdown
# 🛡️ Weekend Warrior Backend

Бэкенд-часть проекта для организации активного отдыха. Построен на **Go (Golang)** с использованием фреймворка **Gin** и ORM **GORM**.

---

## 🚀 Быстрый старт

### 1. Системные требования
Убедитесь, что у вас установлены:
* **Go 1.21+**
* **Git**

### 2. Установка зависимостей
Склонируйте репозиторий и скачайте необходимые модули:
```bash
go mod tidy
```

### 3. Запуск сервера
Выполните команду из корня проекта:
```bash
go run cmd/server/main.go
```
Сервер будет доступен по адресу: `http://localhost:8080`

---

## 📚 Документация API (Swagger)

Мы используем Swagger для интерактивного тестирования API (как в FastAPI).

### Как открыть:
1. Запустите сервер.
2. Перейдите по ссылке: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

### Как обновить документацию:
Если вы изменили аннотации (комментарии) в коде, документацию нужно перегенерировать:
1. Установите утилиту (один раз):
   ```bash
   go install [github.com/swaggo/swag/cmd/swag@latest](https://github.com/swaggo/swag/cmd/swag@latest)
   ```
2. Сгенерируйте файлы:
   ```bash
   swag init -g cmd/server/main.go
   ```

---

## 📂 Структура проекта

* `cmd/server/main.go` — Точка входа, настройка роутов и запуск.
* `internal/handlers/` — Логика обработки запросов (контроллеры).
* `internal/models/` — Структуры данных (GORM модели).
* `internal/db/` — Подключение к базе данных (SQLite).
* `docs/` — Автогенерируемые файлы Swagger.

---

## 🛠️ Стек технологий
* **Language:** Go (Golang)
* **Framework:** [Gin Gonic](https://gin-gonic.com/)
* **Database:** SQLite (файл `warriors.db`)
* **ORM:** [GORM](https://gorm.io/)

---

## 🤝 Правила работы (Git Workflow)

1. **Main Branch:** Прямой пуш в `main` запрещен.
2. **Features:** Для каждой задачи создаем ветку: `git checkout -b feature/your-task-name`.
3. **PR:** Изменения вносятся через **Pull Request** с обязательным код-ревью.
4. **Ports:** Если порт `8080` занят после остановки, используйте:
   ```bash
   sudo fuser -k 8080/tcp
   ```
```
