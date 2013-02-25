package main

import (
	"io/ioutil"
	"log"
	"fmt"
	"net/http"
	"html/template"
	"strings"
	"strconv"
	"time"
)

type Metric struct {
	UnixTime int64
	Year int
	Month int
	Day int
	Hour int
	Minute int
	Second int
	Value float64
}

func buildMetric (m Metric, e []string) Metric {
	m.Year, _ = strconv.Atoi(e[0])
	m.Month, _ = strconv.Atoi(e[1]) // JS month starts from 0 fixme
	m.Day, _ = strconv.Atoi(e[2])
	m.Hour, _ = strconv.Atoi(e[3])
	m.Minute, _ = strconv.Atoi(e[4])
	m.Second, _ = strconv.Atoi(e[5])
	return m
}

func parse_time(t string) string {
	tm, err := strconv.Atoi(t)
	if err != nil {
		log.Println("Atoi failed: ", err)
	}
	ptime := time.Unix(int64(tm), 0).Format("2006,1,2,15,04,05")
	return ptime
}

func getMetrics(query string) []Metric {
	toGraph := []Metric{}

	url := fmt.Sprintf("http://localhost:4242/q?%s", query) // should be a parameter
	resp, err := http.Get(url)
	if err != nil {
		log.Print("GET display: ", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	strBody := string(body)
	metrics := strings.Split(strBody, "\n")

	for i := range metrics {
		if len(metrics[i]) > 2 {
			m := Metric{}
			metric := strings.Split(metrics[i], " ")
			t := parse_time(metric[1])
			m = buildMetric(m, strings.Split(t,","))
			m.Value, _ = strconv.ParseFloat(metric[2], 64)
			toGraph = append(toGraph, m)
		}
	}
	return toGraph
}

func main(){

	displayHandler := func(w http.ResponseWriter, req* http.Request){
		get := req.URL.Query()
		req.Body.Close()

		url := "?"
		for k,v := range get {
			url = fmt.Sprintf("%s&%s=%s", url, k, v[0] )
		}
		toGraph := getMetrics(url)
		t, _ := template.ParseFiles("metric.html")
		t.Execute(w, toGraph)
	}

	http.HandleFunc("/display", displayHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}
