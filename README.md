# Lottery

Сервис для проведения лотерей.

## Запуск

### Локальный запуск
> пример конфига в [config.local.yaml](config.local.yaml)
```bash
go run ./app/server -config ./config.yaml
```

### Запуск через Docker

```bash
docker compose -f compose-local.yaml up --build
```

## Архитектура

Проект построен на чистой архитектуре с разделением на слои:

- `internal/` - основной код приложения
  - `auth/` - аутентификация и авторизация
  - `draw/` - управление тиражами
  - `export/` - экспорт данных
  - `lottery/` - бизнес-логика лотерей
  - `payment/` - платежи
  - `result/` - результаты тиражей
  - `ticket/` - управление билетами
  - `config/` - конфигурация
  - `server/` - HTTP сервер
  - `storage/` - работа с БД

- `pkg/` - общие пакеты
- `app/` - точка входа
- `migrations/` - миграции БД
- `tests/` - нагрузочное тестирование

## Функциональность

### Аутентификация
- Регистрация пользователя
- Вход в систему
- JWT авторизация

### Лотереи
- 5 из 36
- 6 из 45

### Билеты
- Предварительное создание билетов для тиража
- Пользовательский выбор чисел в билете
- Просмотр билетов
- Проверка статуса билета

### Тиражи
- Создание тиражей
- Проведение тиражей
- Отмена тиражей
- Просмотр результатов тиражей

### Платежи
- Оплата билетов
- Проверка статуса платежа

### Администрирование
- Управление пользователями
- Управление тиражами
- Экспорт данных

## Реализованные нефункциональные требования

### Нагрузочное тестирование
- Время ответа API не более 0.5 секунды.
- Поддержка до 1000 запросов в секунду

Для нагрузочного тестирования используется утилита `k6`.
Результаты  
```
k6 run --vus 10 --duration 5s tests/load/load.js

         /\      Grafana   /‾‾/
    /\  /  \     |\  __   /  /
   /  \/    \    | |/ /  /   ‾‾\
  /          \   |   (  |  (‾)  |
 / __________ \  |_|\_\  \_____/

  █ THRESHOLDS

    errors
    ✓ 'rate<0.01' rate=0.00%

    http_req_duration
    ✓ 'p(95)<500' p(95)=10.68ms


  █ TOTAL RESULTS

    checks_total.......................: 25748   5146.634509/s
    checks_succeeded...................: 100.00% 25748 out of 25748
    checks_failed......................: 0.00%   0 out of 25748

    ✓ статус 200
    ✓ время ответа менее 500ms

    CUSTOM
    errors..................................................................: 0.00%  0 out of 12874

    HTTP
    http_req_duration.......................................................: avg=3.77ms min=229µs    med=1.81ms max=89.46ms p(90)=9.5ms  p(95)=10.68ms
      { expected_response:true }............................................: avg=3.77ms min=229µs    med=1.81ms max=89.46ms p(90)=9.5ms  p(95)=10.68ms
    http_req_failed.........................................................: 0.00%  0 out of 12874
    http_reqs...............................................................: 12874  2573.317255/s

    EXECUTION
    iteration_duration......................................................: avg=3.87ms min=267.62µs med=1.91ms max=91.56ms p(90)=9.61ms p(95)=10.78ms
    iterations..............................................................: 12874  2573.317255/s
    vus.....................................................................: 10     min=10         max=10
    vus_max.................................................................: 10     min=10         max=10

    NETWORK
    data_received...........................................................: 1.5 MB 304 kB/s
    data_sent...............................................................: 1.4 MB 273 kB/s




running (05.0s), 00/10 VUs, 12874 complete and 0 interrupted iterations
default ✓ [======================================] 10 VUs  5s
```

### Логироввание
Различные уровни логирования. Debug, Info, Error.
Ротация логов осуществляется на стороне инфраструктуры согласно принципам [12 факторного приложения](https://12factor.net/ru/logs).

## Интерфейсы
Реализован REST API. Формат запросов/ответов: JSON.
База данных: PostgreSQL 17.
Схема нормализована - вторая нормальная форма.
Индексы для ускорения поиска.

Проведено полное ручное тестирование проекта.  
Приложен каталог API вызовов в Insomnia совместимом формате. Файл: [api-insomnia.yaml](api-insomnia.yaml)
