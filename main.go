package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"sort"
	"sync"
	"time"
)

const FILE = "data/m9.txt"
const BATCH_SIZE = 10000

var wg sync.WaitGroup
var mx sync.Mutex

type stats struct {
	min, sum, max int
	count         int
}

func parseInt(s string) int {
	var result int
	var sign int = 1
	var index int
	var point bool = false

	if s[0] == '-' {
		sign = -1
		index++
	}

	for ; index < len(s); index++ {
		ch := s[index]

		if ch != '.' {
			digit := int(ch - '0')
			result = result*10 + digit
		} else {
			point = true
		}
	}

	if !point {
		result *= 10
	}

	return result * sign
}

func findSeparator(s string) int {
	for i, ch := range s {
		if ch == ';' {
			return i
		}
	}
	return -1
}

func processBatch(batch []string, localMap map[string]stats) {
	for _, text := range batch {
		sepIndex := findSeparator(text)
		num := parseInt(text[sepIndex+1:])
		name := text[:sepIndex]

		v, ok := localMap[name]
		if !ok {
			localMap[name] = stats{min: num, max: num, sum: num, count: 1}
		} else {
			if num < v.min {
				v.min = num
			}
			if num > v.max {
				v.max = num
			}
			v.sum += num
			v.count++
			localMap[name] = v
		}
	}
}

func mergeMaps(globalMap, localMap map[string]stats) {
	mx.Lock()
	defer mx.Unlock()
	for k, v := range localMap {
		if globalV, ok := globalMap[k]; ok {
			if v.min < globalV.min {
				globalV.min = v.min
			}
			if v.max > globalV.max {
				globalV.max = v.max
			}
			globalV.count += v.count
			globalV.sum += v.sum
			globalMap[k] = globalV
		} else {
			globalMap[k] = v
		}
	}
}

var batchPool = sync.Pool{
	New: func() interface{} {
		return make([]string, 0, BATCH_SIZE)
	},
}

func main() {
	runtime.GOMAXPROCS(4)
	start := time.Now()

	cpuProfile, err := os.Create("cpu.prof")
	if err != nil {
		fmt.Println("Error creating CPU profile:", err)
		return
	}
	defer cpuProfile.Close()

	if err := pprof.StartCPUProfile(cpuProfile); err != nil {
		fmt.Println("Error starting CPU profile:", err)
		return
	}
	defer pprof.StopCPUProfile()

	f, err := os.Open(FILE)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(bufio.NewReaderSize(f, 16*1024*1024))

	globalMap := make(map[string]stats)
	batch := batchPool.Get().([]string)
	lineCount := 0

	for scanner.Scan() {
		batch = append(batch, scanner.Text())
		lineCount++

		if lineCount == BATCH_SIZE {
			wg.Add(1)
			localMap := make(map[string]stats)
			go func(batch []string) {
				defer wg.Done()
				processBatch(batch, localMap)
				mergeMaps(globalMap, localMap)
				batchPool.Put(batch[:0])
			}(batch)
			lineCount = 0
			batch = batchPool.Get().([]string)
		}
	}
	wg.Wait()

	names := make([]string, 0, len(globalMap))
	for k := range globalMap {
		names = append(names, k)
	}

	sort.Strings(names)

	for _, name := range names {
		v := globalMap[name]
		fmt.Printf("%s:[%.1f %.1f %.1f], ", name, float32(v.min)/10, (float32(v.sum)/float32(v.count))/10, float32(v.max)/10)
	}

	memProfile, err := os.Create("mem.prof")
	if err != nil {
		fmt.Println("Error creating memory profile:", err)
		return
	}
	defer memProfile.Close()

	runtime.GC()
	if err := pprof.WriteHeapProfile(memProfile); err != nil {
		fmt.Println("Error writing memory profile:", err)
	}

	fmt.Println("\nTime elapsed:", time.Since(start), "\nLines remain:", lineCount)
}
