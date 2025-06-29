package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		panic("Error loading .env file")
	}

	servers := []Server{
		NewSimpleServer(os.Getenv("SERVER_1_API_URL")),
		NewSimpleServer(os.Getenv("SERVER_2_API_URL")),
		NewSimpleServer(os.Getenv("SERVER_3_API_URL")),
		NewSimpleServer(os.Getenv("SERVER_4_API_URL")),
		NewSimpleServer(os.Getenv("SERVER_5_API_URL")),
	}

	ConnectRedis()
	lb := NewLoadBalancer("8080", servers)
	
	handleCoin := func(rw http.ResponseWriter, r *http.Request) {
		lb.ServeProxy(rw, r)
	}
	handleChart := func (rw http.ResponseWriter, r *http.Request) {
		lb.ChartServerProxy(rw, r)
	}

	http.HandleFunc("/{name}", handleCoin)
	http.HandleFunc("/{name}/market_chart", handleChart)

	fmt.Printf("serving request at 'localhost:%s'", lb.Port)
	http.ListenAndServe(":"+lb.Port, nil)
}
