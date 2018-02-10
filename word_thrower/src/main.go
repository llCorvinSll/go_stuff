package main

import (
	"net/http"
	"sync"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/", handler)

	http.ListenAndServe("localhost:8000", nil)
}


const WORD_LEN = 26

type WordCounter struct {
	mut       sync.Mutex
	words     []string
	itemsLeft int
}

var initialWords = WordCounter{
	words:[]string{
		"Warty Warthog",
		"Hoary Hedgehog",
		"Breezy Badger",
		"Dapper Drake",
		"Edgy Eft",
		"Feisty Fawn",
		"Gutsy Gibbon",
		"Hardy Heron",
		"Intrepid Ibex",
		"Jaunty Jackalope",
		"Karmic Koala",
		"Lucid Lynx",
		"Maverick Meerkat",
		"Natty Narwhal",
		"Oneiric Ocelot",
		"Precise Pangolin",
		"Quantal Quetzal",
		"Raring Ringtail",
		"Saucy Salamander",
		"Trusty Tahr",
		"Utopic Unicorn",
		"Vivid Vervet",
		"Wily Werewolf",
		"Xenial Xerus",
		"Yakkety Yak",
		"Zesty Zapus",
	},
	itemsLeft:WORD_LEN,
}


func handler(writer http.ResponseWriter, request *http.Request) {
	initialWords.mut.Lock();
	defer initialWords.mut.Unlock()

	if initialWords.itemsLeft <= 0 {
		fmt.Fprint(writer, "no words");

		return
	}

	nextIndex := rand.Intn(initialWords.itemsLeft)

	fmt.Fprintf(writer, "Next is %q", initialWords.words[nextIndex])

	initialWords.words[nextIndex], initialWords.words[initialWords.itemsLeft-1] = initialWords.words[initialWords.itemsLeft-1], initialWords.words[nextIndex]

	initialWords.itemsLeft--
}