# Итоговый проект

## 1. Краткое описание проекта
Проект представляет собой простейший планировщик задач с веб-интерфейсом и бэком на Go, клиент-серверное взаимодействие реализовано посредством REST API, в качестве СУБД выступает SQLite

Сервис позволяет создавать новые задачи с определенными правилами повторения, а также просматривать, редактировать и удалять их

## 2. Реализованные задачи со звездочкой
1. Возможность определять порт по переменной окружения
2. Создание докер образа

## 3. Запуск кода
 - Для запуска сервера локально необходимо в терминале выполнить сборку проекта находясь в папке проекта
    ```sh
        go build
    ```
    а затем запустить проект
    ```sh
        ./go_final_project 
    ```
 - Для взаимодействия с сервером через веб-интерфейс необходимо при запущенном сервере открыть в браузере
http://localhost:7540/

## 4. Тестирование
- Для запуска автотестов необходимо в терминале, находись в папке проекта, выполнить
    ```sh
            go test ./tests
    ```
- [Коллекция методов API для Postman](https://www.postman.com/timur-tikhomirov/workspace/yandex/request/)

## 5. Docker
- Создание образа
    ```sh
            docker build --tag my_app:v1 .
    ```
- Запуск контейнера 
    ```sh
            docker run -p 7540:7540 -d my_app:v1
    ```