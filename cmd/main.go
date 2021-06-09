package main

import (
	"bufio"
	"fmt"
	"net/http"
	"time"
)

type NetSpace struct {
	Success   bool    `json:"success"`
	DayChange float64 `json:"daychange"`
	Netspace  int64   `json:"netspace"`
	Timestamp int64   `json:"timestamp"`
}
type Market struct {
	Success   bool    `json:"success"`
	Price     float64 `json:"price"`
	Daymin    float64 `json:"daymin"`
	Daymax    float64 `json:"daymax"`
	Daychange float64 `json:"daychange"`
	Timestamp int64   `json:"timestamp"`
}

func main() {
	fmt.Println("Hello,Go!")
	netspaceURL := "https://api.chiaprofitability.com/netspace"
	xchURL := "https://api.chiaprofitability.com/market"

	ReqNetspaceByURL(netspaceURL)
	time.Sleep(2 * time.Second)
	ReqMarketByURL(xchURL)
}

func ReqNetspaceByURL(strURL string) {
	resp, err := http.Get(strURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)

	scanner := bufio.NewScanner(resp.Body)

	for i := 0; scanner.Scan() && i < 5; i++ {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func ReqMarketByURL(strURL string) {
	resp, err := http.Get(strURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)

	scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan() && i < 5; i++ {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
