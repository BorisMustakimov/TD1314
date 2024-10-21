# Final
Финальный проект
# Описание
В итоговом задании нужно написать на Golang веб-сервер, который реализует функциональность простейшего планировщика задач.Это задание на проверку и закрепление навыков по написанию веб-сервера, работе с REST API и базами данных.

В планировщике должны храниться задачи, каждая из которых содержит дату дедлайна, а так же заголовок с комментарием.Задачи могут повторяться ежегодно, через какое-то количество дней , в определенные дни месяца или недели.Обычные задачи при выполнении будут удаляться.

API имеет следующие операции:
1. Добавить задачу
2. Получить список 
3. Удалить задачу
4. Получить параметры задачи
5. Изменить параметры задачи
6. Отметить задачу как выполненную

# Тесты
1. go test -run ^TestApp$ ./tests
2. go test -run ^TestDB$ ./tests
3. go test -run ^TestNextDate$ ./tests
4. go test -run ^TestAddTask$ ./tests
5. go test -run ^TestTasks$ ./tests
6. go test -run ^TestEditTask$ ./tests
7.1 go test -run ^TestDone$ ./tests
7.2 go test -run ^TestDelTask$ ./tests
7.3 go test ./tests - Для прохождения всех тестов сразу
