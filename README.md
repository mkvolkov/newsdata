# newsdata

Сервер fiber, оперирующий данными MySQL.

Для работы необходим установленный MySQL с
созданной БД newsdata и таблицами.

Проверка на докеризованной MySQL:

1. Собрать образ БД с тестовыми значениями.

```
make ms_build
```

2. Запустить контейнер.

```
make ms_run
```

3. Запустить сервер. Перед запуском лучше подождать
около 20 секунд после шага 2, пока контейнер не поднялся.

```
make run
```

4. Остановить контейнер MySQL:

```
make ms_stop
```

5. Удалить контейнер и образ

```
make ms_clean
```