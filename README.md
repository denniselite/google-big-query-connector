Google Big Query Service
==============

Сервис отправляет входящие события в Google big query

Пример конфигурации
--------------
```
HTTP-порт для входящих запросов
listen: 3003

rabbit:
  host: 127.0.0.1
  port: 5672
  username: guest
  password: guest

// Параметры подключения к таблице Google Big Query
bigQuery:
  dataSet: dataSet
  projectID: projectID
  tableName: tableName
```