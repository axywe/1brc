package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	f, err := os.Open("data/m8.txt")
	if err != nil {
		fmt.Println(err)
	}
	scanner := bufio.NewScanner(f)

	m := make(map[string][4]float64)
	for scanner.Scan() {
		text := scanner.Text()
		parts := strings.Split(text, ";")
		num, _ := strconv.ParseFloat(parts[1], 64)
		arr, ok := m[parts[0]]
		if !ok {
			m[parts[0]] = [4]float64{num, num, num, 1}
		} else {
			if num < arr[0] {
				arr[0] = num
			}
			arr[1] = math.Ceil((arr[1]+(num-arr[1])/(arr[3]+1))*10) / 10
			if num > arr[2] {
				arr[2] = num
			}
			arr[3] = arr[3] + 1
			m[parts[0]] = arr
		}
	}

	fmt.Println(m, len(m))
	return
}
