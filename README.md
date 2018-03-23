# Запуск сети
_Все действия в каталоге processor_
## Установить зависимости
```
cd processor 
npm i
```
## Удалить все предыдущие контейнеры
```
docker rm -f $(docker ps -aq) && yes | docker network prune
```
## Запустить сеть
```
docker-compose -f network.yaml up
```
## Запустить transaction processor
```
node index.js
```
# Тестовые запросы
_Все действия в каталоге client_
## Установить зависимости
```
cd client 
npm i
```
## Положить данные в блокчейн
```
node index.js
```
## Прочитать данные из блокчейна
```
node check.js
```

# Troubleshooting
### Не проходят транзакции после перезапуска processor-а
#### Шаг 1 Перезапусите сеть - выполните команды в каталоге processor:
```
docker rm -f $(docker ps -aq) && yes | docker network prune`
docker-compose -f network.yaml up
или - запуск в фоне
docker-compose -f network.yaml up -d
```
#### Шаг 2 Запустите processor в каталоге processor:

```
node index.js
```

# Полезные команды

## Удалить сеть 
```

docker rm -f $(docker ps -aq) && yes | docker network prune
```
## Удалить все none контейнеры
```
docker rmi $(docker images -f "dangling=true" -q)
```
## Удалить образы по имени
```
docker rmi $(docker images | awk '$1 ~ /fabric/ { print $3}')
```
## Собрать образ проекта tfa_backend, tfa_frontend, tfa_cabinet
```
docker build -t allatrack/blockchain_tfa_backend .
docker build -t allatrack/blockchain_tfa_frontend .
docker build -t allatrack/blockchain_tfa_cabinet .
```# blockchain-2fa-net
