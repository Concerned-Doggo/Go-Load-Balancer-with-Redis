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
	handleHome := func(rw http.ResponseWriter, r *http.Request) {
		lb.ServeProxy(rw, r)
	}
	handleCoin := func(rw http.ResponseWriter, r *http.Request) {
		lb.ServeProxy(rw, r)
	}

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/coin/{name}", handleCoin)

	fmt.Printf("serving request at 'localhost:%s'", lb.Port)
	http.ListenAndServe(":"+lb.Port, nil)
}
