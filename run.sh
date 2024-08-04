#!/bin/bash

TIMEFORMAT="Время выполнения - %lR"

case "$1" in
    1)
        FILE="main.go"
        ;;
    2)
        FILE="1/main.go"
        ;;
    *)
        echo "Неверный аргумент. Укажите 1 или 2."
        exit 1
        ;;
esac

go build $FILE
time {
    ./main --race
}

rm main