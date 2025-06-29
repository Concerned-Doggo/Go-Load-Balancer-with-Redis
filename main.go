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
		enableCors(&rw, r)
		lb.ServeProxy(rw, r)
	}
	handleChart := func (rw http.ResponseWriter, r *http.Request) {
		enableCors(&rw, r)
		lb.ChartServerProxy(rw, r)
	}

	http.HandleFunc("/{name}", handleCoin)
    http.HandleFunc("/{name}/market_chart", handleChart)


	fmt.Printf("serving request at 'localhost:%s'", lb.Port)
	http.ListenAndServe(":"+lb.Port, nil)
}

func enableCors(w *http.ResponseWriter, r *http.Request) {
    origin := r.Header.Get("Origin")
	if origin == "https://crypto-insight1.netlify.app" || origin == "http://localhost:5173" || origin == "http://localhost:4173" {
        (*w).Header().Set("Access-Control-Allow-Origin", origin)
    }
    (*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
    (*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
}

