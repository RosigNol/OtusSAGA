В этом ДЗ вы научитесь реализовывать распределенную транзакцию.


Описание/Пошаговая инструкция выполнения домашнего задания:
Можно использовать приведенный ниже сценарий для интернет-магазина или придумать свой.
Дефолтный сценарий:
Реализовать сервисы "Платеж", "Склад", "Доставка".
Для сервиса "Заказ", в рамках метода "создание заказа" реализовать механизм распределенной транзакции (на основе Саги или двухфазного коммита).
Во время создания заказа необходимо:

в сервисе "Платеж" убедиться, что платеж прошел
в сервисе "Склад" зарезервировать конкретный товар на складе
в сервисе "Доставка" зарезервировать курьера на конкретный слот времени.
Если хотя бы один из пунктов не получилось сделать, необходимо откатить все остальные изменения.
На выходе должно быть:
описание того, какой паттерн для реализации распределенной транзакции использовался
команда установки приложения (из helm-а или из манифестов). Обязательно указать в каком namespace нужно устанавливать и команду создания namespace, если это важно для сервиса.
тесты в postman
В тестах обязательно
использование домена arch.homework в качестве initial значения {{baseUrl}}

##О проекте: использовался паттерн SAGA

##Установка

kubectl apply -f ./db
helm install postgres -f values.yaml bitnami/postgresql --set volumePermissions.enabled=true
kubectl apply -f ./order-service
kubectl apply -f ./item-service
kubectl apply -f ./payment-service
kubectl apply -f ./dtm
##Тестирование: сгенерированна коллкция postman (SAGA.postman_collection.json)
##Схема работы: приложена в файле SAGA.drawio.png
