#!/bin/bash

TIMEFORMAT="Time elapsed - %lR"

case "$1" in
    1)
        FILE="main.go"
        ;;
    2)
        FILE="1/main.go"
        ;;
    *)
        echo "Invalid argument. Specify 1 or 2."
        exit 1
        ;;
esac

go build $FILE
time {
    ./main --race
}

rm main