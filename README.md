## Описание проекта

## Выполненные задания

- Порт извне сервера через переменную окружения ✅
- Путь к файлу через переменную окружения ✅
- Поддержка w, m ❌
- Выбор задач через поле поиска ✅
- Аутентификация ✅
- Докер образ  ❌
- README ✅

## Инструкция для локального запуска проекта

Для локального запуска проекта нужно указать значения переменных окружения, сделать это можно несколькими способами:

- Указать значения в файле `config.yaml`. Путь к конфигу также можно задать через переменную окружения `CONFIG_PATH`, если ее не задать, то будет использован путь по умолчанию и предполагаться, что конфиг находится в корне проекта
- Указать значения переменных окружения при запуске проекта
``` bash
HOST="host" PORT="port" DB_PATH="db_path" PASSWORD="password" SECRET="secret" go run cmd/scheduler/main.go
```

``` bash
export HOST="host"
export PORT="port"
export DB_PATH="db_path"
export PASSWORD="password"
export SECRET="secret"
go run cmd/scheduler/main.go
```

## Инструкция по запуску тестов

Для запуска тестов нужно ввести `go test ./tests` из корневой дериктории проекта

## Инструкция по сборке и запуску проекта через докер

Сначала нужно собрать докер образ следующей командой:

``` bash
docker build -t scheduler-app .
```

Дальше нужно запустить контейнер на основе созданного докер образа, прописав переменные оркужения при запуске:

``` bash
docker run -p 7540:7540 \
    -e PORT=7540 \
    -e HOST=0.0.0.0 \
    -e DB_PATH=scheduler.db \
    -e AUTH_PASSWORD=1234 \
    -e AUTH_SECRET=aadfs9fhg-9134hf-981h5fg8h12=f9uq=80g1=38g1=39g \
    scheduler-app
```

Либо прописать нужные переменные окружения в Dockerfile 

``` Dockerfile
ENV PORT=7540
ENV HOST=0.0.0.0
ENV DB_PATH=scheduler.db
ENV AUTH_PASSWORD=qewrdsaf
ENV AUTH_SECRET=aadfs9fhg-9134hf-981h5fg8h12=f9uq=80g1=38g1=39g
```

И запустить контейнер командой 

``` bash
docker run -p 7540:7540 scheduler-app
```