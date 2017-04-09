package main

import (
	"os"
	"fmt"
	"log"
	"bufio"
	"hash/fnv"
	"io/ioutil"
	"unicode"
	"strings"
	"time"
)

const APOSTROPHE = 39

var (
	hashMap = make(map[uint32]string)
)
// Returns hash of the string
func getHash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(strings.ToLower(s)))
	return h.Sum32()
}

// Returns true if word is in dictionary, else false
func isInDictionary(word string) bool {
	if _, ok := hashMap[getHash(word)]; ok {
		return true
	}

	return  false
}

// Loads dictionary into map
func loadDictionary(dictionary string, c1 chan int) {
	wordsInDictionary := 0
	// Open file for reading
	file, err := os.Open(dictionary)
	if err != nil {
		log.Fatal(err)
		return
	}
	// Close this file
	defer file.Close()
	// Read every line
	scanner := bufio.NewScanner(file)
	// Hash and add every new value to map
	for scanner.Scan() {
		wordsInDictionary++
		hashMap[getHash(scanner.Text())] = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return
	}

	c1 <- wordsInDictionary
}
// Returns file name which needs to be checked
func getFileName() string {
	fileName := ""
	fmt.Print("Enter file name: ")
	fmt.Scanf("%s", &fileName)

	return fileName
}

func spellCheck(c2 chan int) {
	word := ""
	index := 0
	misspellings := 0
	wordsInText := 0
	data, _ := ioutil.ReadFile(getFileName())
	fmt.Println("\nMISSPELLED WORDS:\n")
	// Look for each byte
	for _,byte := range data {
		// If current character is letter or apostrophe - add it to string
		if unicode.IsLetter(rune(byte)) || (byte == APOSTROPHE && index > 0) {
			word += string(byte)
			index++
			// If current character is digit - ignore
		} else if unicode.IsDigit(rune(byte)) {
			continue
			// Else we've got a whole word
		} else if index > 0 {
			wordsInText++
			// If this word is not in the dictionary - display it
			if !isInDictionary(word) {
				fmt.Printf("%s\n", word)
				misspellings++
			}

			word = ""
			index = 0
		}
	}

	c2 <- misspellings
	c2 <- wordsInText
}

func showResults(c1 chan int, c2 chan int) {
	fmt.Println("\nWORDS IN DICTIONARY: ", <- c1)
	fmt.Println("WORDS MISSPELLED: ", <- c2)
	fmt.Println("WORDS IN TEXT: ", <- c2)
}

func main() {
	var ch1 chan int = make(chan int)
	var ch2 chan int = make(chan int)

	go loadDictionary("large", ch1)
	go spellCheck(ch2)

	time.Sleep(10 * time.Second)
	showResults(ch1, ch2)
}