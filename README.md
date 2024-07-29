# Файлы для итогового задания

В директории `api` находятся обработчики реализованного API веб-сервера.

В директории `cmd` находится файл main.go.

В директории `config` происходит загрузка переменных окружения из файла .env и инициализация структуры конфигурации.

В директории `docs` находится пользовательская документация Swagger.

В директории `middleware` находится функционал авторизации.

В директории `model` находятся общие для всего приложения константы и описание общих объектов.

В директории `repository` находится весь функционал работы с БД.

В директории `tests` находятся тесты для проверки API, которое должно быть реализовано в веб-сервере.

В директории `usecases` находится бизнес логика приложения.

Директория `web` содержит файлы фронтенда.

Приложение поделено на независимые слои: Транспорт, Бизнес-логика, Работа с БД. Тем самым реализована чистая архитектура.

В файле `.env` задаются значения переменных окружений:
- `TODO_PORT` - порт. Пример "7540", "8080".
- `TODO_DBFILE` - относительный или абсолютный путь к файлу БД. Пример "./scheduler.db".
- `TODO_PASSWORD` - пароль для авторизации. Пример "12345".
- `TODO_LOGLEVEL` - уровень логирования. Пример "INFO", "DEBUG".