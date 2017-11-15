package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	var r = flag.Uint64("r", 10, "rate: how many api calls should executes in one second")
	var f = flag.String("f", "domain.txt", "file: specify textfile written domain name")
	var d = flag.Int64("d", 10, "duration: how many secods should continues")
	flag.Parse()
	sleep_time := time.Duration(1e9 / uint64(*r))
	file := string(*f)
	duration := time.Duration(*d) * time.Second

	if f == nil {
		fmt.Println("Error: must specify some file for -f")
		os.Exit(1)
	}

	var fp *os.File
	var err error

	fp, err = os.Open(file)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	var list []string
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		list = append(list, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	for _, domain := range list {
		go autodig(domain, sleep_time)
	}
	time.Sleep(duration)
	fmt.Println("exit status 0")
}

func autodig(domain string, sleep_time time.Duration) {
	var count int
	for {
		_, err := net.LookupHost(domain)
		if err != nil {
			fmt.Println(err)
			count++
			if count >= 3 {
				fmt.Println("exit status 1")
				os.Exit(1)
			}
			continue
		}
		time.Sleep(sleep_time)
		count = 0
	}
}
