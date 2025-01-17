## gRPC сервис на Go для получения приложения курса USDT

## Требования.

- **Go**
- **Docker**: Docker используется для деплоя сервисов с помощью `docker-compose`.
- **Компилятор Protobuf (protoc)**: Необходим для генерации Go-кода из `.proto` файлов.

## Настройка окружения.

Проект использует переменные окружения для конфигурации различных сервисов. Шаблонный файл переменных окружения предоставлен как `.env.example`. 
Перед запуском проекта необходимо создать файл `.env` и заполнить его реальными значениями.

## Отредактируйте файл `.env`** и замените значения переменных на реальные:

   ```env
    LOG_LEVEL=debug   # Уровень вывода логов: info,debug,error 
    GRPC_PORT=50051  # Порт для gRPC сервера
    HTTP_HOST=exchange_rate:8080 #Адрес и порт HTTP сервера
    HTTP_PORT=8080 #Порта http сервера
    PROM_PORT=9090 #Порт доступа Prometheus c localhost
    PG_HOST=postgres #Host PostgreSQL
    PG_PORT=5432 #Порт PostgreSQL
    PG_USER=template_service  # Имя пользователя PostgreSQL
    PG_PASSWORD=template_service  # Имя пользователя PostgreSQL
    PG_DATABASE=template_service # Имя базы данных PostgreSQL
    JAEGER_PORT=16686 #Порт доступа Jaeger с localhost
   ```

## Структура проекта

Проект организован по следующим директориям:

###  `cmd/service/`

Содержит точку входа в приложение — файл `main.go`. Здесь определяется основной запуск сервиса.

- **`main.go`**: Главный файл, отвечающий за запуск и конфигурацию приложения.

###  `config/`

Содержит файлы конфигурации проекта.

- **`config.go`**: Здесь определяется структура конфигурационных данных и логика их загрузки.

###  `db/`

Содержит файлы, связанные с базой данных.

- **`migrations/`**: Содержит SQL файлы миграций базы данных. 
- 
###  `deploy/`

Содержит файлы, связанные с деплоем приложения с помощью Docker.

- **`docker-compose/`**:
    - **`docker-compose.yaml`**: Конфигурация для запуска сервисов через Docker Compose.
    - **`Dockerfile`**: Файл для сборки Docker-образа приложения.

###  `internal/`

Содержит основные пакеты сервиса, которые не должны быть доступны извне.

- **`application`**:
  **Слой приложения**: 

    - grpc/ - Реализация gRPC сервера
    - dto/ -  Обработка DTO

- **`models` && `modules`**:
  **Бизнес-логика и сущности **: 

    - models/ - Определение доменных сущностей.
    - modules/ - Реализация бизнес-логики и сценариев взаимодействия с БД

- **`infrastructure`**:
  **Слой инфраструктуры**: Компоненты, которые взаимодействуют с внешними системами

    - router/ - Определение роутера для съема метрик Prometheus и трейсинга с помощью Jaeger
    - apiGarantex/ - Взаимодействие с внешним API биржи

- **`generated`**:
  **Содержит автоматически сгенерированный код.**

###  `pkg/`

**Содержит вспомогательные пакеты для работы с БД и логгером**

### `schemas/proto/`

Содержит исходные файлы `.proto` для генерации gRPC и Protobuf кода.

### Makefile

- **build - сборка исполняемого файла в корне проекта*
- **test - запуск тестов, написаны юнит тесты для проверки работоспособности логики приложения**
- **lint - запуск линтера golangci-lint**
- **docker-build - сборка образа приложения(имя образа(exchagerate) задается командой)** 
- **run - запуск всех контейнеров в фоновов режиме**
- **stop - остановка всех контейнеров в фоновом режиме**

### Запуск приложения

   1. git clone
   2. Копирование файла .env.example и внесение актуальных переменных окружения
   3. Сборка образа приложения make docker-build
   4. Запуск всех контейнеров make run
      4.1. Конфигурация базы через файл ./deploy/postgresql/postgresql.conf 
      или напрямую в строке command в файле ./deploy/docker-compose/docker-compose.yaml
   5. Остановка контейнеров через make stop 
      

