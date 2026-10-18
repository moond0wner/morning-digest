# 🌅 morning-digest

Личный Telegram-бот с утренним дайджестом. Присылает всё самое нужное в одном сообщении: погоду, состояние инвестиционного портфеля и свежие статьи. Написан на Go.

## ✨ Функционал

- [x] **Дайджест** — каждый день в заданное время присылает дайджест по всем фичам.
- [x] **Погода** — текущая температура и ощущения для вашего города.
- [x] **Инвестиции** — отслеживание портфеля в реальном времени через Tinkoff Invest API.
- [x] **Хабр** — топовые статьи с RSS-ленты Хабра.

## 🛠 Стек технологий

- **Язык:** Go (1.22+)
- **Telegram:** [go-telegram/bot](https://github.com/go-telegram/bot)
- **Конфигурация:** YAML + переменные окружения
- **Логирование:** `zap`
- **Контейнеризация:** Docker & Docker Compose
- **Архитектура:** Feature-based / Clean Architecture

## 🚀 Быстрый старт

### Требования

- Go 1.22 или выше
- Docker & Docker Compose (опционально)
- Токен Telegram-бота (получить у [@BotFather](https://t.me/BotFather))
- API-ключ OpenWeather
- Токен Tinkoff Invest API

### Установка и запуск

1. **Клонируйте репозиторий:**
   ```bash
   git clone https://github.com/moond0wner/morning-digest.git
   cd morning-digest
   ```
2. **Настройте переменные окружения:**
   ```
    BOT_TOKEN=
    OPEN_WEATHER=
    INVEST_TOKEN=
    INVEST_ACCOUNT_ID=
    TELEGRAM_ACCOUNT_ID=
    
    LOGGER_LEVEL=DEBUG
    LOGGER_FOLDER=out/logs
    ```
3. **Запуск:**
   
   *Локально:*
   ```bash
   make run
   ```

   *Docker:*
   ```bash
   make deploy
   ```

### Структура проекта
```
├── cmd/                  # Точки входа (main.go, Dockerfile)
├── internal/             # Приватный код приложения
│   ├── core/             # Ядро (конфиг, логгер, шедулер, ретраи, телеграм)
│   └── features/         # Бизнес-логика (погода, инвестиции, хабр, дайджест)
├── out/                  # Директория для вывода (логи, база данных)
├── .env                  # Переменные окружения
├── config.yaml           # Конфигурация приложения 
├── docker-compose.yml    # Конфигурация Docker Compose
├── Dockerfile            # Многоступенчатая сборка Docker
└── Makefile              # Удобные команды (run, deploy, lint)
```
