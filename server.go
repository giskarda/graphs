package main

import (
	"io/ioutil"
	"log"
	"fmt"
	"net/http"
	"strings"
)


func GetMetrics(query string) map[string]string {
	toGraph := map[string]string {}
	url := fmt.Sprintf("http://localhost:4242/q?%s", query)
	resp, err := http.Get(url)
	if err != nil {
		log.Print("GET display: ", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	strBody := string(body)
	metrics := strings.Split(strBody, "\n")
	for m := range metrics {
		if len(metrics[m]) > 2 {
			metric := strings.Split(metrics[m], " ")
			toGraph[metric[1]] = metric[2]
		}
	}
	defer resp.Body.Close()
	return toGraph
}

func main(){

	display := func(w http.ResponseWriter, req* http.Request){
		get := req.URL.Query()
		req.Body.Close()
		url := "?"

		for k,v := range get {
			url = fmt.Sprintf("%s&%s=%s", url, k, v[0] )
		}
		toGraph := GetMetrics(url)
		for k,v := range toGraph {
			fmt.Println(k,v)
		}
	}

	http.HandleFunc("/display", display)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}
