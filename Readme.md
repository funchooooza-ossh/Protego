# ADR-001: Архитектура инфраструктурного Auth-сервиса

## Контекст

В рамках MSA-стенда необходимо реализовать инфраструктурный сервис `auth-backend`, который будет выполнять централизованную аутентификацию и авторизацию всех запросов к API уровня бэкенда, включая внутренние сервисы, фронт и административные панели.

Сервис должен быть:
- независимым и масштабируемым
- безопасным (включая поддержку JWT, blacklist, блокировки)
- гибким в расширении (RBAC, аудит, federation)

## Решение

Принята следующая архитектура:

### 1. Ядро авторизации и ролей — PostgreSQL
- users
- roles
- user-role binding
- permissions (по необходимости)
- блокировки (временные и постоянные)
- refresh sessions (fingerprint, TTL)

### 2. Session Store — Redis
- refresh blacklist (отозванные токены)
- access blacklist (временные отключения)
- rate-limit на логин/refresh по IP/email
- TTL для single-use токенов

### 3. Audit Log — ClickHouse
- logins (успешные/неудачные)
- refresh события
- логауты
- попытки заблокированных пользователей

## Технологии

| Слой              | Технология              | Назначение                        |
|-------------------|--------------------------|-----------------------------------|
| Web transport     | FastAPI + Msgspec        | REST API с OpenAPI-доками         |
| JWT               | python-jose / pyjwt      | HS256 токены                      |
| DB ORM            | SQLAlchemy 2.0 (async)   | Модели и транзакции               |
| Redis client      | redis.asyncio            | Сессии, блокировки                |
| ClickHouse client | clickhouse-driver        | Аудит                             |
| DI                | кастомный контейнер      | Тестируемость                     |
| Test stack        | pytest + httpx           | Юниты + интеграционные тесты      |
| CI/CD             | Makefile + Docker + K8s  | Автоматизация + деплой            |

## Структура проекта

```
auth_backend/
├── application/
│   ├── services/
│   │   └── authenticator.py
│   └── usecases/
│       ├── login.py
│       ├── refresh.py
│       └── logout.py
├── domain/
│   ├── models/
│   │   └── user.py
│   ├── contracts/
│   │   └── authenticator.py
│   └── policies/
│       └── blocking.py
├── infrastructure/
│   ├── db/
│   │   ├── models/
│   │   ├── repositories/
│   │   └── migrations/
│   ├── redis/
│   │   └── blacklist_store.py
│   └── logging/
│       └── clickhouse_writer.py
├── api/
│   ├── routes/
│   │   ├── login.py
│   │   └── refresh.py
│   └── deps.py
├── composition/
│   ├── containers/
│   ├── di.py
│   └── create_app.py
├── tests/
│   ├── unit/
│   └── integration/
├── k8s/
│   ├── base/
│   └── overlays/
│       └── local/
├── Makefile
├── pyproject.toml
└── README.md
```

## Последствия

- Сервис можно использовать как:
  - точку входа для пользователей и админов
  - endpoint для `forwardAuth` в Traefik
- Все токены централизованно валидируются
- Возможность расширения до SSO (например, OAuth2 provider)

## Далее

- [ ] Определить схему БД (DDL)
- [ ] Подготовить фабрики UoW + DI
- [ ] Реализовать контракты `Authenticator`, `SessionStore`, `AuditLogWriter`
- [ ] Поднять окружение в `k8s/overlays/local`
