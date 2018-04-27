
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
&& docker build . -f sawtooth-dev-go -t sawtooth-dev-go \ 
&& cd ..

``` 
#### 2 Шаг - Клонировать репозиторий sawtooth-core
Для сборки проекта нам необходимо иметь sawtooth-core который содержит все sdk для работы с блокчейном
Выполняется из корня проекта
```
git clone --branch v1.0.2 https://github.com/hyperledger/sawtooth-core.git sawtooth-core
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
# залить в dockerhub
docker push 5478545378/sawtooth-tfa-s-tp-go
docker push 5478545378/sawtooth-tfa-sc-tp-go
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
## Остановить все контейнеры
```
docker stop $(docker ps -a -q)
```
## Удалить все none контейнеры
```
docker rmi $(docker images -f "dangling=true" -q)
```
## Удалить образы по имени
```
docker rmi $(docker images | awk '$1 ~ /fabric/ { print $3}')
```
## Запуск сети
```
rm -rf networks/config/conf.d/* && docker-compose -f networks/network-dev.yaml up
```
## Запуск в фоне обработчиков транзакций - только для DEBUG-га
```
nohup go run go/src/tfa/service/main.go --connect=tcp://172.18.0.2:4004 --family=tfa --version=0.1 --verbose  > /dev/null 2>&1 &
nohup go run go/src/tfa/service_client/main.go --connect=tcp://172.18.0.2:4004 --family=kaztel --version=0.1 --verbose  > /dev/null 2>&1 &
```
## Сборка proto файлов для go программы
```
protoc --go_out=handler service_client.proto 
protoc --go_out=handler service.proto 
```
