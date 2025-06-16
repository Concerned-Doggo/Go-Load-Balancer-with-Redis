package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Server interface {
	Address() string
	IsAlive() bool
	Serve(rw http.ResponseWriter, r *http.Request)
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
	coin := r.URL.Query().Get("name")
	// if data is present in our cache then send the data
	redisData := GetRedisData(coin)
	if redisData != "" {
		fmt.Println("sent from cache")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(redisData))
		return
	}

	response, err := http.Get(s.Address())
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
