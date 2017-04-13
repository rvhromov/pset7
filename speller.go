package main

import (
	"os"
	"fmt"
	"log"
	"bufio"
	"io/ioutil"
	"unicode"
	"strings"
)

const APOSTROPHE = 39

var (
	dictionary = make(map[string]string)
	wordsInDictionary = 0
	misspellings = 0
	wordsInText = 0
)

// Returns true if word is in dictionary, else false
func isInDictionary(word string, channel chan bool) {
	_, success := dictionary[word];
	channel <- success
}

// Loads dictionary into map
func loadDictionary() {
	file, err := os.Open("large")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	// Read every line
	scanner := bufio.NewScanner(file)
	// Add every new value to map
	for scanner.Scan() {
		wordsInDictionary++
		dictionary[scanner.Text()] = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// Returns text from file
func readTextFromFile() []byte {
	text, _ := ioutil.ReadFile("austinpowers.txt")
	return text
}

// Shows words which aren't in dictionary
func showMisspelledWords(words []string) {
	channel := make(chan bool)

	for _,word := range words {
		go isInDictionary(word, channel)

		if ! <- channel {
			fmt.Printf("%s\n", word)
			misspellings++
		}
	}
}

// Returns an array which contains only words
func getWordsFromText(text []byte) []string {
	word := ""
	var  words []string

	for _,byte := range text {
		// If current character is letter or apostrophe - add it to string
		if unicode.IsLetter(rune(byte)) || (byte == APOSTROPHE && len(word) > 0) {
			word += strings.ToLower(string(byte))
			// Else we got a whole word
		} else if len(word) > 0 {
			words = append(words, word)
			wordsInText++
			word = ""
		}
	}

	return words
}

func showResults() {
	fmt.Println("\nWORDS IN DICTIONARY: ", wordsInDictionary)
	fmt.Println("WORDS MISSPELLED: ", misspellings)
	fmt.Println("WORDS IN TEXT: ", wordsInText)
}

func main() {
	loadDictionary()
	showMisspelledWords(getWordsFromText(readTextFromFile()))
	showResults()
}