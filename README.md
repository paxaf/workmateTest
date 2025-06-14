# workmateTest - Управление цитатами

[![Go Version](https://img.shields.io/badge/go-1.23%2B-blue.svg)](https://golang.org/)

Cервис для управления задачами с in-memory хранилищем, реализованный на чистой архитектуре и имитацией I/O bound нагрузки

## Особенности
- Проект написан только на stdlib
- Создание, получение, обновление и удаление задач
- In-memory хранилище реализрванное для многопоточного доступа
- Чистая архитектура с разделением слоёв
- RESTful API
- Юнит-тесты для обработчиков
- Рантайм шедулер для иммитации нагрузки

## Структура проекта (Clean Architecture)
├── cmd # Точка входа  
├── internal  
│ ├── app # Инициализация приложения  
│ ├── controller # Логика обработчиков  
│ │  ├── middleware #Логика роутинга  
│ ├── entity # Бизнес-сущности (Task)  
│ ├── repository # Интерфейсы хранилища  
│ │ ├── engine # In-memory реализация  
│ ├── usecase  # Интерфейсы и реализация бизнес-логики  
└── go.sum  # файл для корректной сборки  
└── build.log # проверка сборки с запуском тестов с флагом -race  
└── docker-compose.yml # Запуск сервиса в контейнерной среде  
└── Dockerfile # Файл для сборки образа сервиса  
└── README.md # Описание проекта

## Требования
- Go 1.23+
- Для тестов: `go test` и `gcc` компилятор для использования `-race`

## Установка и запуск
```bash
# Клонировать репозиторий
git clone https://github.com/paxaf/workmateTest.git
cd workmateTest

# Запустить сервер
go run cmd/server/main.go

# Или запустить docker-compose
docker-compose up -d --build
## (не запускайте их одновременно, они слушают один и тот же порт)
```
## API Endpoints

## API Endpoints

| Метод   | Путь           | Описание                 |
|---------|----------------|--------------------------|
| POST    | `/tasks`      | Создать задачу         |
| GET     | `/tasks`      | Получить все задачи      |
| GET     | `/tasks?id=`      | Получить задачу по её id  |
| DELETE  | `/tasks/{id}`  | Удалить задачу           |

## Запуск тестов
Если установлен `gcc` в корне проекта можно использовать команду в `bash`
```bash
go test ./... -race -v
```
Если по каким то причинам `gcc` у вас в нет или флаг `-race` не работает корректно на вашей локальной машине. Тогда для теста с флагом `-race` можно открыть `Dockerfile` изменить 
```Dockerfile 
CGO_ENABLED=1 # Устанавливаем значение на 1
...
RUN go test ./... -race -cover -v # убираем `#` перед этой строчкой
```
Запускаем команду 
```bash
docker-compose build | tee build.log
```

И смотрим как проходят тесты с флагом `-race`.

Перед билдом бинарника для использования в контейнере ``обязательно`` возвращаем всё в исходное состояние, иначе проект не запустится.

