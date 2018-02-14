package main

import (
	"net/http"
	"fmt"
	"math/rand"
	"time"
	"os"
	"log"
	"bufio"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	strings := reader()

	results := randomGenereator(strings, 3, 4)


	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, <-results);
	})

	http.ListenAndServe("localhost:8000", nil)
}



const NO_WORDS = "NO WORDS LEFT"




func reader() <-chan string {
	out := make(chan string, 1)

	file, err := os.Open("words.txt")
	if err != nil {
		log.Fatal("CANT open file")
	}

	scanner := bufio.NewScanner(file)

	go func() {
		for {
			if scanner.Scan() {
				out <- scanner.Text()
			} else {
				out <- NO_WORDS
			}
		}
	}()

	return  out
}

func worker(results chan<- string, initialArray []string, strReader <- chan string, workerId int) {
	for {
		actualPoolLength := len(initialArray);

		if actualPoolLength == 0 {
			results <- "Sorry!"
			continue
		}

		newIndex := rand.Intn(actualPoolLength)

		res := initialArray[newIndex]

		newItem := <-strReader

		if newItem != NO_WORDS {
			initialArray[newIndex] = newItem
		} else {
			initialArray = append(initialArray[:newIndex], initialArray[newIndex+1:]...)
		}

		fmt.Println("--------------------------------")
		fmt.Printf("worker number %v \n", workerId)
		fmt.Printf("%v \n", initialArray)
		fmt.Println("--------------------------------")

		results <- res
	}
}


func randomGenereator(strings <-chan string, poolSize int, concurrency int) <-chan string {
	out := make(chan string, 1)

	for i := 0; i < concurrency; i++ {
		var s []string

		for j := 0; j < poolSize; j++ {
			s = append(s, <- strings)
		}

		go worker(out, s, strings, i)
	}


	return out
}