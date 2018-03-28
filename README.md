
# Building Transaction Families 

## Инструкция для разработчика 
### Сборка Transaction Families
Bin файлы собраны и находятся в каталоге 
```
go/bin
```
Влучае изменения кода transaction family выполните слудующие шаги:

#### 1 Шаг - Сборка docker image - сборщик Go программ
```
# Build:
cd docker \
&& docker build . -f sawtooth-dev-go -t sawtooth-dev-go\ 
&& cd ..

``` 
#### 2 Шаг - Клонировать репозиторий sawtooth-core
Для сборки проекта нам необходимо иметь sawtooth-core который содержит все sdk для работы с блокчейном
Выполняется из корня проекта
```
git clone --branch v1.0.1 git@github.com:hyperledger/sawtooth-core.git sawtooth-core
```
#### 3 Шаг - Сборка проекта Go программы
Выполняется из корня проекта
```
cp $(pwd)/scripts/build_go_bin $(pwd)/sawtooth-core/bin
docker run -v $(pwd)/sawtooth-core:/project/sawtooth-core \
           -v $(pwd)/go:/project/tfa/go \
           sawtooth-dev-go
```
#### 3 Шаг - Сборка docker image - Transaction Families которые будут запукаться
Выполняется из корня проекта
```
./scripts/build_go_images
``` 

## Запуск сети - 1 валидатор(для разработки)
```

```


# Полезные команды
## Golang переменные окружения 
```
export GOROOT=/usr/local/go
export GOPATH=$HOME/go:<путь к проекту>blockchain-2fa-net/go
```
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