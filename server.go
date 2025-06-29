package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

type Server interface {
	Address() string
	IsAlive() bool
	Serve(rw http.ResponseWriter, r *http.Request)
	ServeChart(rw http.ResponseWriter, r *http.Request)
}

type SimpleServer struct {
	address string
	Proxy   *httputil.ReverseProxy
}

func (s *SimpleServer) Address() string {
	return s.address
}

func (s *SimpleServer) IsAlive() bool {
	return true
}

func (s *SimpleServer) Serve(w http.ResponseWriter, r *http.Request) {

	// fmt.Println(strings.Split(r.URL.Path, "/"))
	coin := strings.Split(r.URL.Path, "/")[1]
	extraParam := ""
	// if data is present in our cache then send the data
	redisData := GetRedisData(coin)
	if redisData != "" {
		fmt.Println("sent from cache")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(redisData))
		return
	}
	
	url := s.GetParamsEncodedURL(coin, "usd", extraParam)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add(os.Getenv("API_HEADER"), os.Getenv("API_KEY"))

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer response.Body.Close()

	res, err := io.ReadAll(response.Body)
	fmt.Printf("%T\n", res)

	// set the data in our cache
	SetRedisData(coin, res)

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (s *SimpleServer) ServeChart(w http.ResponseWriter, r *http.Request){
	urlParams := strings.Split(r.URL.Path, "/")
	coin := urlParams[1]
	chart := urlParams[2]

	redisKey := string(coin + chart)
	fmt.Println("redisKey: ", redisKey);

	redisChartData := GetRedisData(redisKey)
	if redisChartData != "" {
		fmt.Println("sent chart data from cache")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(redisChartData))
		return
	}

	url := s.GetParamsEncodedURL(coin, "usd", chart)

	fmt.Println("char url: " , url)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add(os.Getenv("API_HEADER"), os.Getenv("API_KEY"))

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer response.Body.Close()

	res, err := io.ReadAll(response.Body)
	fmt.Printf("%T\n", res)

	// set the data in our cache
	SetRedisData(redisKey, res)

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

}

func (s *SimpleServer) GetParamsEncodedURL(coin string, currency string, extraParam string) string {
	baseUrl := s.Address()
	baseUrl += coin
	if(extraParam != ""){
		baseUrl += "/" + extraParam
	}

	endPoint, err := url.Parse(baseUrl)
	if err != nil{
		fmt.Println("error parsing baseurl: ", endPoint, err)
		return "";
	}

	// common query params
	queryParams := url.Values{}
	queryParams.Set("vs_currency", currency)

	if(extraParam != ""){
		// they are asking for the chart
		queryParams.Set("days", "90")
		queryParams.Set("interval", "daily")
	}
	
	endPoint.RawQuery = queryParams.Encode()

	return endPoint.String()
}


func NewSimpleServer(address string) *SimpleServer {
	serverURL, err := url.Parse(address)
	if err != nil {
		panic(fmt.Sprintf("Error parsing server URL %s: %v", address, err))
	}

	return &SimpleServer{
		address: address,
		Proxy:   httputil.NewSingleHostReverseProxy(serverURL),
	}
}
