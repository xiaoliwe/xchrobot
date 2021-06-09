package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type NetSpace struct {
	Success   bool    `json:"success"`
	DayChange float64 `json:"daychange"`
	Netspace  int64   `json:"netspace"`
	Timestamp string  `json:"timestamp"`
}
type Market struct {
	Success   bool    `json:"success"`
	Price     float64 `json:"price"`
	Daymin    float64 `json:"daymin"`
	Daymax    float64 `json:"daymax"`
	Daychange float64 `json:"daychange"`
	Timestamp string  `json:"timestamp"`
}
type XCH struct {
	Netspace  string
	Price     string
	Daychange string
}
type Notify struct {
	Msgtype  string `json:"msgtype"`
	Markdown struct {
		Content string `json:"content"`
	} `json:"markdown"`
}

var httpClient = &http.Client{Timeout: 10 * time.Second}

func main() {
	ns := new(NetSpace)
	market := new(Market)

	netspaceURL := "https://api.chiaprofitability.com/netspace"
	xchURL := "https://api.chiaprofitability.com/market"

	getJson(netspaceURL, ns)
	fmt.Println(ns.Netspace)
	fmt.Println(ns.DayChange)

	getJson(xchURL, market)
	fmt.Println(market.Price)

	//Post
	xch := new(XCH)
	xch.Daychange = fmt.Sprintf("%f", ns.DayChange)
	xch.Netspace = fmt.Sprintf("%d", ns.Netspace)
	xch.Price = fmt.Sprintf("%f", market.Price)

	robotURL := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=68d83069-cde9-493f-9081-34537f132084"
	postXCH(robotURL, xch)
}

func getJson(url string, target interface{}) error {
	r, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func postXCH(url string, xchInfo *XCH) string {

	contents := fmt.Sprintf("今日币价简报:\n >全网算力: <font color=\"info\">%s EiB</font>\n>新增算力: <font color=\"comment\">%s PiB</font>\n>当前币价: <font color=\"warning\">%s USD</font>",
		xchInfo.Netspace, xchInfo.Daychange, xchInfo.Price)

	notify := new(Notify)
	notify.Msgtype = "markdown"
	notify.Markdown.Content = contents

	fmt.Println(notify)

	postBody, _ := json.Marshal(notify)
	respBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(url, "application/json", respBody)
	if err != nil {
		log.Fatalf("Post failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	log.Printf(sb)
	return sb
}
