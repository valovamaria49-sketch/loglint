# loglint

Линтер для проверки лог-записей в Go (log/slog, go.uber.org/zap).

## Установка

```bash
go install ./cmd/loglint@latest
```

## Запуск

```bash
go run ./cmd/loglint ./...
```

## Интеграция с golangci-lint

Сборка плагина (из корня проекта):

```bash
go build -buildmode=plugin -o ./bin/loglint.so ./analysis
```

В `.golangci.yml`:

```yaml
linters:
  settings:
    custom:
      loglint:
        path: ./bin/loglint.so
        description: "Линтер для проверки лог-записей"
```

Запуск:

```bash
golangci-lint run
```

## Конфигурация

Файл `loglint_config.json` в рабочей директории (рядом с go.mod):

```json
{
  "sensitive_keywords": ["password", "token", "api_key", "secret"],
  "enable_lowercase_rule": true,
  "enable_ascii_rule": true,
  "enable_english_rule": true,
  "enable_sensitive_rule": true
}
```

| Параметр | Описание |
|----------|----------|
| `sensitive_keywords` | Список слов, при наличии которых в логе — ошибка |
| `enable_lowercase_rule` | Проверка: сообщение начинается со строчной буквы |
| `enable_ascii_rule` | Проверка: только ASCII-буквы, цифры, пробел (без спецсимволов и эмодзи) |
| `enable_english_rule` | Проверка: только английские буквы |
| `enable_sensitive_rule` | Проверка: отсутствие чувствительных данных по ключевым словам |

### Применение изменений в конфиге

После изменения `loglint_config.json` нужно **перезапустить** `golangci-lint` (если используется интеграция):

```bash
# Сбросить кеш и перезапустить
golangci-lint cache clean
golangci-lint run
```

## Правила

1. Лог-сообщения должны начинаться со строчной буквы.
2. Только английский язык и ASCII (буквы, цифры, пробел).
3. Без спецсимволов и эмодзи.
4. Без потенциально чувствительных данных (password, token и т.п.).

## Документация

### Локальный просмотр

```bash
# Установите godoc (однократно)
go install golang.org/x/tools/cmd/godoc@latest

# Посмотрите информацию
go doc -all ./analysis
```
