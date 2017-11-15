package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	var fp *os.File
	var err error

	if len(os.Args) < 2 {
		fp = os.Stdin
	} else {
		fp, err = os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
		defer fp.Close()
	}

	var list []string
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		list = append(list, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	for _, domain := range list {
		go dig(domain)
	}
	time.Sleep(60 * time.Second)
	fmt.Println("exit status 1")
}

func dig(domain string) {
	var count int
	for {
		_, err := net.LookupHost(domain)
		if err != nil {
			fmt.Println(err)
			count++
			if count >= 3 {
				os.Exit(1)
			}
			continue
		}
		time.Sleep(1000000)
		count = 0
	}
}
