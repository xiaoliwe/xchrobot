package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/robfig/cron/v3"
)

type NetSpace struct {
	Success   bool        `json:"success"`
	DayChange float64     `json:"daychange"`
	Netspace  json.Number `json:"netspace"`
	Timestamp int64       `json:"timestamp"`
}
type Market struct {
	Success   bool    `json:"success"`
	Price     float64 `json:"price"`
	Daymin    float64 `json:"daymin"`
	Daymax    float64 `json:"daymax"`
	Daychange float64 `json:"daychange"`
	Timestamp int64   `json:"timestamp"`
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
	crontab := cron.New(cron.WithLogger(cron.VerbosePrintfLogger(log.New(os.Stdout, "cron:", log.LstdFlags))))
	_, err := crontab.AddFunc("CRON_TZ=Asia/Shanghai 0 19 * * *", func() { handlerPost() })
	if err != nil {
		log.Printf("main-->cron error:%s", err)
		return
	}
	crontab.Start()
	defer crontab.Stop()

	select {}
	//handlerPost()
}
func handlerPost() {
	ns := new(NetSpace)
	market := new(Market)

	namespaceURL := "https://api.chiaprofitability.com/netspace"
	xchURL := "https://api.chiaprofitability.com/market"

	err := getJson(namespaceURL, ns)
	if err != nil {
		return
	}
	err2 := getJson(xchURL, market)
	if err2 != nil {
		return
	}

	f, _ := ns.Netspace.Float64()
	eb := f / 1152921504606846976

	AllPower := fmt.Sprint(eb)[0:5]
	XCHPrice := fmt.Sprint(market.Price)[0:6]
	UpdatedTime := time.Unix(ns.Timestamp, 0)
	NewPower := fmt.Sprint(((ns.DayChange / 100.00) * eb) * 1024)[0:6]

	fmt.Printf("Success is: %v\n", ns.Success)
	fmt.Printf("Netspace(EiB) is: %v\n", AllPower)
	fmt.Printf("Daychange is: %v\n", ns.DayChange)
	fmt.Printf("NewChange(PiB) is: %v\n", NewPower)
	fmt.Printf("Timestamp is: %v\n", UpdatedTime)
	fmt.Printf("XCH Price(USD): %s", XCHPrice)

	//Post
	xch := new(XCH)
	xch.Daychange = fmt.Sprintf("%v", NewPower)
	xch.Netspace = fmt.Sprintf("%v", AllPower)
	xch.Price = fmt.Sprintf("%v", XCHPrice)

	robotURL := "YOUR WEBHOOK ADDRESS OF WECOM"
	postXCH(robotURL, xch)
	fmt.Println("Push XCHPrice Success!")
}

func getJson(url string, target interface{}) error {
	r, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(r.Body)

	return json.NewDecoder(r.Body).Decode(target)
}

func postXCH(url string, xchInfo *XCH) string {

	contents := fmt.Sprintf("??????????????????:\n >????????????: <font color=\"info\">%s EiB</font>\n>????????????: <font color=\"comment\">%s PiB</font>\n>????????????: <font color=\"warning\">%s USD</font>\n????????????:<font color=\"warning\"> 0.0003 XCH/TiB</font>",
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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	log.Printf(sb)
	return sb
}
