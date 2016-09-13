package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func parseWords(fname string) (words [][]string, err error) {
	f, err := os.Open(fname)
	if err != nil {
		return
	}
	defer func() {
		err = f.Close()
		return
	}()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		words = append(words, strings.Fields(scanner.Text()))
	}
	return
}

func cleanWords(words [][]string) [][]string {
	const requiredWords = 1<<12 + 1<<4
	if len(words) < requiredWords {
		log.Fatalf("insufficient words: need %d, only have %d", requiredWords, len(words))
	}
	return words[:requiredWords]
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		log.Fatal("no file provided")
	}
	fname := flag.Arg(0)
	words, err := parseWords(fname)
	if err != nil {
		log.Fatal(err)
	}
	words = cleanWords(words)

	fmt.Println("package wordenc")
	fmt.Println("")
	fmt.Println("var wordList = [...][]string{")

	for _, wordGroup := range words {
		fmt.Printf("\t")
		fmt.Printf("{")
		wordStrs := make([]string, len(wordGroup))
		for i, word := range wordGroup {
			wordStrs[i] = fmt.Sprintf("\"%s\"", word)
		}
		fmt.Print(strings.Join(wordStrs, ", "))
		fmt.Printf("},\n")
	}

	fmt.Println("}")
}
