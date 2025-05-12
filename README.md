# Веб приложение "Производственная гимнастика" (monorepo)

Веб приложение для управления сервисом "Производственная гимнастика", позволяет настраивать расписание, отслеживать статистику, и управлять разминочными комплексами для клиентов приложения.

## Клиентское приложение

https://github.com/RO1T/gym-client

## Stack

Server: Golang

Frontend: Vue/Quasar

Protocol: GRPc

## Запуск веб приложения

Запуск системы + генерация протоген файлов.

make bundle-gen

Запуск системы без генерации

make bundle

## Доступные пути

http://localhost:8083/ - docs

http://localhost:8085/ - envoy

http://localhost:9090/ - server

http://localhost:8080/ - frontend (основная точка входа)


## Ошибки

### ./entrypoint error 

dos2unix gym-server/entrypoint.sh

git update-index --chmod=+x gym-server/entrypoint.sh
