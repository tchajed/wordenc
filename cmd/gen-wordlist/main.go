package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

func parseWords(fname string) (words []string, err error) {
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
		words = append(words, scanner.Text())
	}
	return
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

	fmt.Println("package wordenc")
	fmt.Println("")
	fmt.Println("var wordList = [...]string{")

	for _, word := range words {
		fmt.Printf("\t\"%s\",\n", word)
	}

	fmt.Println("}")
}
