package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"time"
)

var (
	concurrencyFactor int
	host              = flag.String("host", "localhost:8000", "address to listen")
	poolSize          = flag.Int("pool", 3, "number of items in each random pool")
	sourceFile        = flag.String("source", "words.txt", "file with items")
	drainMessage      = flag.String("drain_message", "Sorry!!!", "answer then no values left")
)

func main() {
	flag.Parse()
	concurrencyFactor = runtime.NumCPU()
	rand.Seed(time.Now().UnixNano())

	log.Printf("listening address [%v]", *host)
	log.Printf("number of pools [%v]", concurrencyFactor)
	log.Printf("poolsize [%v]", *poolSize)

	strings := reader(sourceFile)

	results := randomGenerator(strings, *poolSize, concurrencyFactor)

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, <-results)
	})

	log.Fatal(http.ListenAndServe(*host, nil))
}

func reader(path *string) <-chan string {
	out := make(chan string, 1)

	file, err := os.OpenFile(*path, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatalf("can't open file [%v]", *path)
	} else {
		log.Printf("source file opened successful [%v]", *path)
	}

	scanner := bufio.NewScanner(file)

	go func() {
		defer func() {
			err := file.Close()
			if err != nil {
				log.Fatalf("fail to close file [%v] /n", file.Name())
			} else {
				log.Printf("file closed successfuly [%v] \n", file.Name())
			}

			close(out)
		}()

		for {
			if scanner.Scan() {
				out <- scanner.Text()
			} else {
				return
			}
		}
	}()

	return out
}

func worker(results chan<- string, initialArray []string, strReader <-chan string, workerId int) {
	for {
		actualPoolLength := len(initialArray)

		if actualPoolLength == 0 {
			results <- *drainMessage
			continue
		}

		newIndex := rand.Intn(actualPoolLength)

		res := initialArray[newIndex]

		newItem, err := <-strReader

		if err {
			initialArray[newIndex] = newItem
		} else {
			initialArray = append(initialArray[:newIndex], initialArray[newIndex+1:]...)
		}

		log.Printf("worker[%v] cap - %v \n", workerId, len(initialArray))

		results <- res
	}
}

func randomGenerator(strings <-chan string, poolSize int, concurrency int) <-chan string {
	out := make(chan string, 1)

	for i := 0; i < concurrency; i++ {
		var s []string

		for j := 0; j < poolSize; j++ {
			s = append(s, <-strings)
		}

		go worker(out, s, strings, i)
	}

	return out
}
