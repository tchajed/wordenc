package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		log.Fatal("no file provided")
	}
	fname := flag.Arg(0)
	f, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Println("package wordenc")
	fmt.Println("")
	fmt.Println("var wordList = [...]string{")

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fmt.Printf("\t\"%s\",\n", scanner.Text())
	}

	fmt.Println("}")
}
