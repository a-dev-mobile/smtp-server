FROM golang:1.21.4

# Устанавливаем рабочий каталог в корне контейнера
WORKDIR /go/src/app

# Копируем файлы go.mod и go.sum и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект в контейнер
COPY . .

# Переходим в каталог cmd для сборки приложения
WORKDIR /go/src/app/cmd
RUN go build -o main .

# Экспонируем необходимые порты
EXPOSE 50051 465

# Запускаем скомпилированный бинарник
CMD ["./main"]
