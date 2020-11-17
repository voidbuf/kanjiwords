package main

import (
	"bufio"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var (
	words = []wordEntry{}
)

type wordEntry struct {
	Kanji string
	Freq  int
}

type wordSource interface {
	populateMap()
}

type wikipediaSource struct{}

func (s wikipediaSource) populateMap() {
	file, err := os.Open("sources/wikipedia.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currFreq := 1
	for scanner.Scan() {
		words = append(words, wordEntry{scanner.Text(), currFreq})
		currFreq++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func isJapaneseWord(w string) bool {
	matched, err := regexp.MatchString(`[\p{Han}|\p{Hiragana}|\p{Katakana}]+`, w)
	if err != nil {
		log.Fatal(err)
	}
	return matched
}

type leedsSource struct{}

func (s leedsSource) populateMap() {
	file, err := os.Open("sources/leeds.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currFreq := 1
	notJapanese := 0
	for scanner.Scan() {
		text := scanner.Text()
		split := strings.Split(text, " ")
		word := split[2]
		if isJapaneseWord(word) {
			words = append(words, wordEntry{word, currFreq})
			currFreq++
		} else {
			notJapanese++
		}
	}
	log.Printf("not japanese: %d\n", notJapanese)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func searchKanji(kanji string) []wordEntry {
	matches := []wordEntry{}
	if kanji == "" {
		return matches
	}
	for _, word := range words {
		if strings.Contains(word.Kanji, kanji) {
			matches = append(matches, word)
		}
	}
	return matches
}

func initData(dataSource wordSource) {
	dataSource.populateMap()
}

func main() {
	initData(leedsSource{})
	http.HandleFunc("/favicon.ico", handleFavicon)
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/search", handleSearch)
	log.Println("server listening...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	log.Println("received request at /")
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(200)
	err = t.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}

type searchData struct {
	KanjiSlice []wordEntry
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	kanjiQuery, ok := r.URL.Query()["kanji"]
	kanji := ""
	if ok && len(kanjiQuery[0]) >= 1 {
		kanji = kanjiQuery[0]
	}
	log.Println("received request at /search with: " + kanji)
	matches := searchKanji(kanji)
	funcMap := template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
	}
	t := template.New("search.html").Funcs(funcMap)
	t, err := t.ParseFiles("templates/search.html")
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(200)
	err = t.Execute(w, searchData{matches})
	if err != nil {
		log.Fatal(err)
	}
}

func handleFavicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "favicon.ico")
}
