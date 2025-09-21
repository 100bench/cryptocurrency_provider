# Cryptocurrency Provider

REST API для получения курсов криптовалют. Сбор данных с внешних источников, сохранение в PostgreSQL, агрегаты по времени.

## Запуск

```bash
# Запуск всех сервисов (PostgreSQL, Kafka, приложение)
docker-compose up -d

# Проверка статуса
docker-compose ps
```

Приложение будет доступно на `http://localhost:8080`

## API

- `GET /rates/latest?currencies=BTC,ETH` - последние курсы для указанных валют
- `GET /rates/{currency}/min` - минимальный курс для валюты
- `GET /rates/{currency}/max` - максимальный курс для валюты  
- `GET /rates/{currency}/avg` - средний курс для валюты

### Примеры

```bash
# Получить последние курсы BTC и ETH
curl "http://localhost:8080/rates/latest?currencies=BTC,ETH"

# Получить минимальный курс Bitcoin
curl "http://localhost:8080/rates/BTC/min"
```

**Swagger документация**: `http://localhost:8080/swagger/index.html`

## Технологии

- **Go** - основной язык
- **PostgreSQL** - хранение курсов
- **Kafka** - асинхронная обработка
- **Docker** - контейнеризация
- **CoinDesk API** - источник данных

## Конфигурация

Настройки в `deployment/config/config.yaml`. По умолчанию курсы обновляются каждые 10 секунд.
