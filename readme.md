# SMTP Server для отправки электронной почты

Этот проект представляет собой гибкий SMTP сервер, предназначенный для отправки электронных писем. Он поддерживает работу через gRPC.

## Особенности

- Поддержка множественных SMTP провайдеров.
- Автоматический перебор почтовых серверов для успешной отправки писем.
- Hастраиваемая система логирования.
- Конфигурируемый gRPC сервер для интеграции с другими сервисами.
- Образ Docker доступен на Docker Hub (`wayofdt/smtp-server:latest`).

## Структура проекта

┣ 📂cmd
┃ ┣ 📜main.go
┣ 📂config
┃ ┣ 📜config.example.yaml
┣ 📂internal
┃ ┣ 📂config
┃ ┃ ┗ 📜config.go
┃ ┣ 📂handlers
┃ ┃ ┗ 📂send
┃ ┃ ┃ ┗ 📜send.go
┃ ┗ 📂utils
┃ ┃ ┣ 📜smtp_config.go
┃ ┃ ┗ 📜validation.go
┣ 📂lib
┃ ┗ 📂logger
┃ ┃ ┗ 📂sl
┃ ┃ ┃ ┗ 📜sl.go
┣ 📂proto
┃ ┣ 📜email-sender.pb.go
┃ ┣ 📜email-sender.proto
┃ ┗ 📜email-sender_grpc.pb.go
┣ 📜.gitignore
┣ 📜Dockerfile
┣ 📜go.mod
┗ 📜go.sum


## Конфигурация

Перед запуском необходимо заполнить `config.yaml`, используя `config.example.yaml` в качестве шаблона.

## Запуск

Для запуска сервера выполните команду:
```cmd
go run cmd/main.go
```

## Использование Docker

Запустите образ с помощью следующей команды:
```cmd
docker pull wayofdt/smtp-server:latest
docker run wayofdt/smtp-server:latest
```