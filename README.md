## 🏗️ Структура проекта

```bash
analytics/
│
├── cmd/                        # Точки входа в приложение
│   └── server/
│       └── main.go            # Основной исполняемый файл
│
├── internal/                   # Внутренние пакеты приложения
│   ├── api/                    # API слой
│   │   ├── handlers/           # Обработчики запросов
│   │   ├── middleware/         # Промежуточное ПО
│   │   ├── gen/                # Сгенерированный ogen код
│   │   └── server.go           # Инициализация HTTP сервера
│   │
│   ├── config/                 # Конфигурация приложения
│   │   └── config.go
│   │
│   ├── domain/                 # Доменный слой
│   │   ├── models/             # Доменные модели
│   │   │   └── event.go
│   │   └── services/           # Бизнес-логика
│   │       └── analytics.go
│   │
│   ├── infrastructure/         # Инфраструктурный слой
│   │   ├── kafka/              # Работа с Kafka
│   │   │   ├── consumer.go
│   │   │   └── config.go
│   │   ├── storage/            # Работа с БД
│   │   │   ├── clickhouse/
│   │   │   │   ├── client.go
│   │   │   │   └── repository.go
│   │   │   ├── migrations/      # Миграции для ClickHouse
│   │   │   │   ├── manager.go   # Менеджер миграций
│   │   │   │   └── runner.go    # Запуск миграций
│   │   │   └── repository.go   # Интерфейс репозитория
│   │   └── metrics/            # Работа с метриками
│   │       ├── prometheus.go
│   │       └── metrics.go
│   │
│   └── workers/                # Фоновые рабочие процессы
│       ├── processor.go        # Обработчик сообщений из Kafka
│       └── exporter.go         # Экспорт метрик
│
├── pkg/                        # Публичные пакеты, полезные утилиты
│   ├── logging/                # Логирование
│   └── utils/                  # Утилиты
│
├── api/                        # OpenAPI спецификации
│   └── openapi.yaml
│
├── deployments/                # Инструменты развертывания
│   ├── docker/
│   │   ├── Dockerfile
│   │   └── docker-compose.yaml
│   └── kubernetes/
│       ├── deployment.yaml
│       └── service.yaml
│
├── migrations/                 # SQL миграции
│   ├── clickhouse/             # Миграции для ClickHouse
│   │   ├── 001_initial_schema.sql
│   │   ├── 002_add_indexes.sql
│   │   └── ...
│   └── README.md               # Документация по миграциям
│
├── scripts/                    # Вспомогательные скрипты
│   ├── generate.sh             # Скрипт для генерации кода из OpenAPI
│   └── migrate.sh              # Скрипт для запуска миграций
│
├── go.mod                      # Зависимости Go модуля
├── go.sum
├── .gitignore
├── README.md
└── Makefile                    # Автоматизация сборки
```