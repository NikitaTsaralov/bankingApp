# Система транзакций

## ТЗ:

Реализовать систему транзакций. Она должна:
* принимать запросы на ввод или вывод денег и заносить их в базу данных. 
* отдавать данные по транзакции из бд при запросе инфы по какой-либо транзакции.

## Задания

1. Необходимо реализовать http api для внешних потребителей, которые и будут создавать транзакции и просить информацию о них. Как именно будет выглядеть публичный интерфейс - ваше дело. 

> Оцениваться он будет по удобству, понятности, предсказуемости и прочим вещам, которые мы все любим.

2. Система должна быть доступной и отказоустойчивой, способной жить под высокой нагрузкой - запросы на создание транзакций должны всегда выполнятся рано или поздно, то есть если что-то в сервисе упадёт из-за какого-либо nil pointerа, клиенты не должны терять способность создавать заявки на создание транзакций. Как это реализовать - ваше дело. Рекомендуем использовать брокер сообщений. 

> Ваше решение будет оцениваться по критериям понятности и поддерживаемости - гипотетические члены вашей команды не должны тратить лишнее время на расшифровку вашего замысла, исследуя вашу архитектуру или читая ваш код.

### Обязательно:
* Необходимо использовать postgresql в качестве бд.
* Если у вас будет несколько сервисов, которые общаются друг с другом, они _должны_ общаться по grpc.
* Если будете использовать брокер сообщений, рекомендуем rabbitmq - с ним намного легче начать работать, чем с кафкой.

### Дополнительно: 
* Плюсом будет openapi описание вашего публичного api.
* Также плюсом будет простой compose-файл для всей инфры

## Образец

В качестве абстракции использовалась идея с банковской картой. 

Операции пользователя:

### Аккаунт:
* авторизация с JWT
* просмотр личной информации
* просмотр информации карты (аккаунта)

### История:
* просмотр истории переводов
* просмотр информации конкретного перевода

### Денежные операции:
* пополнение
* снятие наличных

## Коллекция postman

> Коллекция запросов доступна в [файле](./API.postman_collection.json)

## Референсы:

* https://github.com/AleksK1NG/Go-Clean-Architecture-REST-API
* https://github.com/AleksK1NG/Go-gRPC-RabbitMQ-microservice
* https://github.com/antsla/saga
* https://github.com/Duomly/go-bank-backend
* https://github.com/isayme/go-amqp-reconnect/blob/master/rabbitmq/rabbitmq.go

#### Testing manual helpers:

```
pip install pgcli
pgcli -h 127.0.0.1 -p 5432 -U postgres -W postgres -d auth_db
```